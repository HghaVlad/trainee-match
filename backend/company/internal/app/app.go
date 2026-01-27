package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	get_vacancy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
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

	compRepo := repository.NewCompanyRepository(compDB)
	vacRepo := repository.NewVacancyRepo(compDB)

	txManager := infra_postgres.NewTxManager(compDB)

	compGetByIDUc := get_company.NewGetByIDUsecase(compRepo)
	compListUc := list_companies.NewUsecase(compRepo)
	compCreateUc := create_company.NewUsecase(compRepo, txManager)
	compUpdateUc := update_company.NewUsecase(compRepo, txManager)
	compDeleteUc := delete_company.NewUsecase(compRepo, txManager)

	vacGetByIDUc := get_vacancy.NewUsecase(vacRepo)

	companyHandler := handlers.NewProfileHandler(
		compGetByIDUc,
		compCreateUc,
		compListUc,
		compUpdateUc,
		compDeleteUc,
	)

	vacancyHandler := handlers.NewVacancyHandler(
		vacGetByIDUc,
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
