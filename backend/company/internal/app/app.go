package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/indexsearch"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/elastic"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres/repository"
	appredis "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/redis"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/kafka"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/schemaregistry"
	httpapp "github.com/HghaVlad/trainee-match/backend/company/internal/transport/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/handlers"
	compmiddleware "github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/dlq"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/eventhandler"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
	createcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	getcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	listcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	listcompmy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/listmy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/memget"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/moderationstatus"
	removecomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/remove"
	updatecomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	addmember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	listmember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/list"
	membme "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/me"
	removemember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove"
	updatemember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	createvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	getvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/getpublished"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listcompsearch"
	listvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
	vmoderationstatus "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/moderationstatus"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
	removevac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/remove"
	updatevac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type App struct {
	cfg         *config.Config
	HTTPSrv     *http.Server
	outboxRelay *outbox.Relay
	pgDB        *pgxpool.Pool
	rediss      redis.UniversalClient
	kConsumer   *kafka.Consumer
	kProducer   *kafka.Producer
	logger      *slog.Logger
}

//nolint:funlen // app wiring
func Build(ctx context.Context, cfg *config.Config, lgr *slog.Logger) (*App, error) {
	pgDB, err := postgres.ConnectPgxPoolWithLogger(ctx, cfg.Postgres, lgr)
	if err != nil {
		return nil, err
	}

	rediss, err := appredis.NewClient(cfg.Redis)
	if err != nil {
		return nil, err
	}

	schemaRegCl := schemaregistry.NewClient(cfg.SchemaRegistry)
	schemaLocalReg, err := schemaregistry.NewLocalRegistry(ctx, schemaRegCl)
	if err != nil {
		return nil, err
	}

	schemaDecoder := schemaregistry.NewDecoder(schemaLocalReg)
	schemaEncoder, err := schemaregistry.NewEncoder(schemaLocalReg)
	if err != nil {
		return nil, err
	}

	kprClient, err := kafka.NewClientForProducer(cfg.Kafka)
	if err != nil {
		return nil, err
	}
	kProducer := kafka.NewProducer(cfg.Kafka, kprClient, lgr)

	elasticCl, err := elastic.NewClient(cfg.Elastic)
	if err != nil {
		return nil, err
	}

	if err := elastic.Init(ctx, elasticCl); err != nil {
		return nil, err
	}

	compRepo := repository.NewCompanyRepository(pgDB)
	vacRepo := repository.NewVacancyRepo(pgDB)
	memRepo := repository.NewCompanyMemberRepo(pgDB)
	hrProjRepo := repository.NewHrProjectionRepo(pgDB)
	outboxRepo := repository.NewOutboxRepo(pgDB)
	txManager := postgres.NewTxManager(pgDB)

	searchVacRepo := elastic.NewVacancyRepo(elasticCl)

	compCache := appredis.NewRepo[uuid.UUID, company.Company](rediss, "company", lgr)
	vacCache := appredis.NewRepo[uuid.UUID, vacancy.Vacancy](rediss, "vacancy", lgr)
	publicVacCache := appredis.NewRepo[uuid.UUID, getpublished.Response](rediss, "vacancy:public", lgr)
	compListCache := appredis.NewRepo[string, listcomp.Response](rediss, "companies:list", lgr)
	vacListCache := appredis.NewRepo[string, listvac.Response](rediss, "vacancies:list", lgr)
	vacByCompListCache := appredis.NewRepo[string, listcompsearch.Response](rediss, "vacancies_by_comp:list", lgr)

	outboxWriter := outbox.NewWriter(cfg.Outbox, outboxRepo, schemaEncoder)
	outboxRelay := outbox.NewRelay(kProducer, outboxRepo, txManager, cfg.Outbox, lgr)
	dlqSender := dlq.NewSender(cfg.Kafka, kProducer, schemaEncoder)

	compGetByIDUc := getcomp.NewGetByIDUsecase(compRepo, compCache)
	compGetByMem := memget.NewUsecase(compRepo)
	compListUc := listcomp.NewUsecase(compRepo, compListCache)
	compListMy := listcompmy.NewUsecase(compListUc)
	compCreateUc := createcomp.NewUsecase(compRepo, memRepo, outboxWriter, txManager)
	compUpdateUc := updatecomp.NewUsecase(compRepo, memRepo, outboxWriter, txManager, compCache)
	compDeleteUc := removecomp.NewUsecase(compRepo, memRepo, outboxWriter, txManager, compCache)

	compMeUc := membme.NewUsecase(memRepo)
	compAddMemUc := addmember.NewUsecase(memRepo, hrProjRepo, outboxWriter, txManager)
	compListMemUc := listmember.NewUsecase(memRepo)
	compDeleteMemberUc := removemember.NewUsecase(memRepo, outboxWriter, txManager)
	compUpdateMemberUc := updatemember.NewUsecase(memRepo)
	compUpdModerationUc := moderationstatus.NewUsecase(compRepo, outboxWriter, txManager, compCache)

	vacGetByIDUc := getvac.NewUsecase(vacRepo, vacCache, memRepo)
	vacGetPublishedByIDUc := getpublished.NewUsecase(vacRepo, publicVacCache)
	vacList := listvac.NewUsecase(vacRepo, vacListCache)
	searchVac := listvac.NewUsecase(searchVacRepo, vacListCache)
	vacListByComp := listcompsearch.NewUsecase(vacRepo, compRepo, memRepo, vacByCompListCache)
	vacListCompSearch := listcompsearch.NewUsecase(searchVacRepo, compRepo, memRepo, vacByCompListCache)
	vacCreate := createvac.NewUsecase(vacRepo, compRepo, outboxWriter, txManager)
	vacUpdate := updatevac.NewUsecase(vacRepo, compRepo, outboxWriter, vacCache, txManager)
	vacPublish := publish.NewUsecase(vacRepo, compRepo, memRepo, outboxWriter, txManager, vacCache, compCache)
	vacArchive := archive.NewUsecase(
		vacRepo,
		compRepo,
		memRepo,
		outboxWriter,
		txManager,
		vacCache,
		publicVacCache,
		compCache,
	)
	vacUpdModStatus := vmoderationstatus.NewUsecase(
		vacRepo,
		compRepo,
		outboxWriter,
		txManager,
		vacCache,
		publicVacCache,
		compCache,
	)
	vacDelete := removevac.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, publicVacCache, compCache)

	userHrCreate := userhr.NewCreatedUsecase(hrProjRepo)
	searchIndexer := indexsearch.NewUsecase(vacRepo, searchVacRepo)
	eventHandler := eventhandler.NewHandler(
		cfg.KafkaHandling,
		schemaDecoder,
		dlqSender,
		userHrCreate,
		searchIndexer,
		lgr,
	)

	kConsumer, err := kafka.NewConsumer(cfg.Kafka, eventHandler, lgr)
	if err != nil {
		return nil, err
	}

	companyHandler := handlers.NewCompanyHandler(
		compGetByIDUc,
		compGetByMem,
		compCreateUc,
		compListUc,
		compListMy,
		compUpdModerationUc,
		compUpdateUc,
		compDeleteUc,
	)
	memberHandler := handlers.NewMemberHandler(
		compMeUc,
		compAddMemUc,
		compListMemUc,
		compUpdateMemberUc,
		compDeleteMemberUc,
	)
	vacancyHandler := handlers.NewVacancyHandler(
		vacGetByIDUc,
		vacGetPublishedByIDUc,
		vacList,
		searchVac,
		vacListByComp,
		vacListCompSearch,
		vacCreate,
		vacUpdate,
		vacPublish,
		vacArchive,
		vacUpdModStatus,
		vacDelete,
	)

	authMiddleware, err := compmiddleware.NewAuthMiddleware(ctx, cfg)
	if err != nil {
		return nil, err
	}

	routerDeps := &httpapp.RouterDeps{
		CompanyHandler: companyHandler,
		MemberHandler:  memberHandler,
		VacancyHandler: vacancyHandler,
		AuthMiddleware: authMiddleware,
		Logger:         lgr,
	}

	httpRouter := httpapp.NewRouter(routerDeps)

	httpServer := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      httpRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		HTTPSrv:     httpServer,
		pgDB:        pgDB,
		rediss:      rediss,
		kProducer:   kProducer,
		kConsumer:   kConsumer,
		outboxRelay: outboxRelay,
		cfg:         cfg,
		logger:      lgr,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		app.logger.Info("http server starting", "addr", app.HTTPSrv.Addr)

		err := app.HTTPSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http listen and serve: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		app.outboxRelay.Run(ctx)
		return nil
	})

	g.Go(func() error {
		app.kConsumer.Poll(ctx)
		return nil
	})

	return g.Wait()
}

func (app *App) Shutdown(shutdownCtx context.Context) {
	if err := app.HTTPSrv.Shutdown(shutdownCtx); err != nil {
		app.logger.WarnContext(shutdownCtx, "http shutdown error", "err", err)
	}

	app.kProducer.Close()
	app.pgDB.Close()

	if err := app.rediss.Close(); err != nil {
		app.logger.WarnContext(shutdownCtx, "redis shutdown error", "err", err)
	}
}
