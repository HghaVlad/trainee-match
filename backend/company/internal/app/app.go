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
	"golang.org/x/sync/errgroup"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/redis"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/kafka"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/schemaregistry"
	httpapp "github.com/HghaVlad/trainee-match/backend/company/internal/transport/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/handlers"
	compmiddleware "github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
	createcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	getcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	listcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	listcompmy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/listmy"
	removecomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/remove"
	updatecomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	addmember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	removemember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove"
	updatemember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	createvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	getvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/getpublished"
	listvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listbycomp"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
	removevac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/remove"
	updatevac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type App struct {
	cfg         *config.Config
	HTTPSrv     *http.Server
	outboxRelay *outbox.Relay
	pgDB        *pgxpool.Pool
	logger      *slog.Logger
}

//nolint:funlen // app wiring
func Build(ctx context.Context, cfg *config.Config, lgr *slog.Logger) (*App, error) {
	pgDB, err := postgres.ConnectPgxPoolWithLogger(ctx, cfg.Postgres, lgr)
	if err != nil {
		return nil, err
	}

	rediss, err := redis.NewClient(cfg.Redis)
	if err != nil {
		return nil, err
	}

	schemaRegCl := schemaregistry.NewClient(cfg.SchemaRegistry)
	schemaLocalReg, err := schemaregistry.NewLocalRegistry(ctx, schemaRegCl)
	if err != nil {
		return nil, err
	}

	schemaEncoder, err := schemaregistry.NewEncoder(schemaLocalReg)
	if err != nil {
		return nil, err
	}

	kprClient, err := kafka.NewClientForProducer(cfg.Kafka)
	if err != nil {
		return nil, err
	}
	kProducer := kafka.NewProducer(kprClient, lgr)

	compRepo := repository.NewCompanyRepository(pgDB)
	vacRepo := repository.NewVacancyRepo(pgDB)
	memRepo := repository.NewCompanyMemberRepo(pgDB)
	outboxRepo := repository.NewOutboxRepo(pgDB)
	txManager := postgres.NewTxManager(pgDB)

	compCache := redis.NewRepo[uuid.UUID, company.Company](rediss, "company", lgr)
	vacCache := redis.NewRepo[uuid.UUID, vacancy.Vacancy](rediss, "vacancy", lgr)
	publicVacCache := redis.NewRepo[uuid.UUID, getpublished.Response](rediss, "vacancy:public", lgr)
	compListCache := redis.NewRepo[string, listcomp.Response](rediss, "companies:list", lgr)
	vacListCache := redis.NewRepo[string, listvac.Response](rediss, "vacancies:list", lgr)
	vacByCompListCache := redis.NewRepo[string, listbycomp.Response](rediss, "vacancies_by_comp:list", lgr)

	outboxWriter := outbox.NewWriter(cfg.Outbox, outboxRepo, schemaEncoder)
	outboxRelay := outbox.NewRelay(kProducer, outboxRepo, txManager, cfg.Outbox, lgr)

	compGetByIDUc := getcomp.NewGetByIDUsecase(compRepo, compCache)
	compListUc := listcomp.NewUsecase(compRepo, compListCache)
	compListMy := listcompmy.NewUsecase(compListUc)
	compCreateUc := createcomp.NewUsecase(compRepo, memRepo, txManager)
	compAddHrUc := addmember.NewUsecase(memRepo, outboxWriter, txManager)
	compDeleteMemberUc := removemember.NewUsecase(memRepo, outboxWriter, txManager)
	compUpdateMemberUc := updatemember.NewUsecase(memRepo)
	compUpdateUc := updatecomp.NewUsecase(compRepo, memRepo, outboxWriter, txManager, compCache)
	compDeleteUc := removecomp.NewUsecase(compRepo, memRepo, outboxWriter, txManager, compCache)

	vacGetByIDUc := getvac.NewUsecase(vacRepo, vacCache, memRepo)
	vacGetPublishedByIDUc := getpublished.NewUsecase(vacRepo, publicVacCache)
	vacList := listvac.NewUsecase(vacRepo, vacListCache)
	vacListByComp := listbycomp.NewUsecase(vacRepo, compRepo, memRepo, vacByCompListCache)
	vacCreate := createvac.NewUsecase(vacRepo, memRepo)
	vacUpdate := updatevac.NewUsecase(vacRepo, memRepo, outboxWriter, vacCache, txManager)
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
	vacDelete := removevac.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, publicVacCache, compCache)

	companyHandler := handlers.NewCompanyHandler(
		compGetByIDUc,
		compCreateUc,
		compListUc,
		compListMy,
		compUpdateUc,
		compDeleteUc,
	)
	memberHandler := handlers.NewMemberHandler(compAddHrUc, compUpdateMemberUc, compDeleteMemberUc)

	vacancyHandler := handlers.NewVacancyHandler(
		vacGetByIDUc,
		vacGetPublishedByIDUc,
		vacList,
		vacListByComp,
		vacCreate,
		vacUpdate,
		vacPublish,
		vacArchive,
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

	return g.Wait()
}

func (app *App) Shutdown(shutdownCtx context.Context) {
	err := app.HTTPSrv.Shutdown(shutdownCtx)
	if err != nil {
		app.logger.WarnContext(shutdownCtx, "shutdown error", "err", err)
	}

	app.pgDB.Close()
}
