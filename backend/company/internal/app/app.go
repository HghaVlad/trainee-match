package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	httpapp "github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
	compmiddleware "github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/redis"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/logger"
	createcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	deletecomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/delete"
	getcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	listcomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	updatecomp "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	addmember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	deletemember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/delete"
	updatemember "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	createvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	deletevac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
	getvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/getpublished"
	listvac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listbycomp"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
	updatevac "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type App struct {
	conf    *config.Config
	HttpSrv *http.Server
	compDB  *sqlx.DB
	logger  *slog.Logger
}

func Build(conf *config.Config) (*App, error) {
	psgConf := postgres.NewConfig(conf)
	lgr := logger.NewSlogLogger()
	compDB, err := postgres.New(psgConf, lgr)
	if err != nil {
		return nil, err
	}

	redisConf := redis.NewConfig(&conf.Redis)
	rediss, err := redis.NewClient(redisConf)
	if err != nil {
		return nil, err
	}

	compRepo := repository.NewCompanyRepository(compDB)
	vacRepo := repository.NewVacancyRepo(compDB)
	memRepo := repository.NewCompanyMemberRepo(compDB)
	txManager := postgres.NewTxManager(compDB)

	compCache := redis.NewRepo[uuid.UUID, company.Company](rediss, "company")
	vacCache := redis.NewRepo[uuid.UUID, vacancy.Vacancy](rediss, "vacancy")
	publicVacCache := redis.NewRepo[uuid.UUID, getpublished.Response](rediss, "vacancy:public")
	compListCache := redis.NewRepo[string, listcomp.Response](rediss, "companies:list")
	vacListCache := redis.NewRepo[string, listvac.Response](rediss, "vacancies:list")
	vacByCompListCache := redis.NewRepo[string, listbycomp.Response](rediss, "vacancies_by_comp:list")

	compGetByIDUc := getcomp.NewGetByIDUsecase(compRepo, compCache)
	compListUc := listcomp.NewUsecase(compRepo, compListCache)
	compCreateUc := createcomp.NewUsecase(compRepo, memRepo, txManager)
	compAddHrUc := addmember.NewUsecase(memRepo)
	compDeleteMemberUc := deletemember.NewUsecase(memRepo)
	compUpdateMemberUc := updatemember.NewUsecase(memRepo)
	compUpdateUc := updatecomp.NewUsecase(compRepo, memRepo, compCache)
	compDeleteUc := deletecomp.NewUsecase(compRepo, memRepo, compCache)

	vacGetByIDUc := getvac.NewUsecase(vacRepo, vacCache, memRepo)
	vacGetPublishedByIDUc := getpublished.NewUsecase(vacRepo, publicVacCache)
	vacList := listvac.NewUsecase(vacRepo, vacListCache)
	vacListByComp := listbycomp.NewUsecase(vacRepo, compRepo, vacByCompListCache)
	vacCreate := createvac.NewUsecase(vacRepo, memRepo)
	vacUpdate := updatevac.NewUsecase(vacRepo, memRepo, vacCache, txManager)
	vacPublish := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)
	vacArchive := archive.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, publicVacCache, compCache)
	vacDelete := deletevac.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, publicVacCache, compCache)

	companyHandler := handlers.NewCompanyHandler(
		compGetByIDUc,
		compCreateUc,
		compListUc,
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

	authMiddleware, err := compmiddleware.NewAuthMiddleware(conf)
	if err != nil {
		return nil, err
	}

	routerDeps := &httpapp.RouterDeps{
		CompanyHandler: companyHandler,
		MemberHandler:  memberHandler,
		VacancyHandler: vacancyHandler,
		AuthMiddleware: authMiddleware,
	}

	httpRouter := httpapp.NewRouter(routerDeps)

	httpServer := &http.Server{
		Addr:         conf.HTTP.Addr,
		Handler:      httpRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		HttpSrv: httpServer,
		compDB:  compDB,
		conf:    conf,
		logger:  lgr,
	}, nil
}

func (app *App) Run() {
	err := app.HttpSrv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("http listening server err", "err", err)
		os.Exit(1)
	}
}

func (app *App) Shutdown(shutdownCtx context.Context) {
	err := app.HttpSrv.Shutdown(shutdownCtx)
	if err != nil {
		app.logger.Warn("shutdown error", "err", err)
	}

	dbErr := app.compDB.Close()
	if dbErr != nil {
		app.logger.Warn("db close error", "err", dbErr)
	}
}
