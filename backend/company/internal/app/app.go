package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/redis"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list_by_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type App struct {
	conf    *config.Config
	httpSrv *http.Server
	compDB  *sqlx.DB
}

func Build(conf *config.Config) (*App, error) {
	psgConf := infra_postgres.NewConfig(conf)
	logger := service_logger.NewSlogLogger()
	compDB, err := infra_postgres.New(psgConf, logger)
	if err != nil {
		return nil, err
	}

	redisConf := infra_redis.NewConfig(&conf.Redis)
	redis, err := infra_redis.NewClient(redisConf)
	if err != nil {
		return nil, err
	}

	compRepo := repository.NewCompanyRepository(compDB)
	vacRepo := repository.NewVacancyRepo(compDB)
	txManager := infra_postgres.NewTxManager(compDB)

	compCache := infra_redis.NewRepo[uuid.UUID, domain.Company](redis, "company")
	vacCache := infra_redis.NewRepo[uuid.UUID, domain.Vacancy](redis, "vacancy")
	compListCache := infra_redis.NewRepo[string, list_companies.Response](redis, "companies:list")
	vacListCache := infra_redis.NewRepo[string, list_vacancy.Response](redis, "vacancies:list")
	vacByCompListCache := infra_redis.NewRepo[string, list_vac_by_comp.Response](redis, "vacancies_by_comp:list")

	compGetByIDUc := get_company.NewGetByIDUsecase(compRepo, compCache)
	compListUc := list_companies.NewUsecase(compRepo, compListCache)
	compCreateUc := create_company.NewUsecase(compRepo)
	compUpdateUc := update_company.NewUsecase(compRepo, compCache)
	compDeleteUc := delete_company.NewUsecase(compRepo, compCache)

	vacGetByIDUc := get_vacancy.NewUsecase(vacRepo, vacCache)
	vacList := list_vacancy.NewUsecase(vacRepo, vacListCache)
	vacListByComp := list_vac_by_comp.NewUsecase(vacRepo, compRepo, vacByCompListCache)
	vacCreate := create_vacancy.NewUsecase(vacRepo, compRepo, txManager)
	vacUpdate := update_vacancy.NewUsecase(vacRepo, vacCache, txManager)
	vacDelete := delete_vacancy.NewUsecase(vacRepo, vacCache)

	companyHandler := handlers.NewCompanyHandler(
		compGetByIDUc,
		compCreateUc,
		compListUc,
		compUpdateUc,
		compDeleteUc,
	)

	vacancyHandler := handlers.NewVacancyHandler(
		vacGetByIDUc,
		vacList,
		vacListByComp,
		vacCreate,
		vacUpdate,
		vacDelete,
	)

	routerDeps := &delivery_http.RouterDeps{
		CompanyHandler: companyHandler,
		VacancyHandler: vacancyHandler,
	}

	httpRouter := delivery_http.NewRouter(routerDeps)

	httpServer := &http.Server{
		Addr:         conf.HTTP.Addr,
		Handler:      httpRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		httpSrv: httpServer,
		compDB:  compDB,
		conf:    conf,
	}, nil
}

func (app *App) Run() {
	err := app.httpSrv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http listening server err: %s\n", err)
	}
}

func (app *App) Shutdown(shutdownCtx context.Context) {
	err := app.httpSrv.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("shutdown error: %v", err)
	}

	dbErr := app.compDB.Close()
	if dbErr != nil {
		log.Printf("db close error: %v", dbErr)
	}
}
