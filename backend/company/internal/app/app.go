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
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
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
	add_member "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	delete_member "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/delete"
	update_member "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
	archive_vacancy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list_by_company"
	publish_vacancy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
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
	memRepo := repository.NewCompanyMemberRepo(compDB)
	txManager := infra_postgres.NewTxManager(compDB)

	compCache := infra_redis.NewRepo[uuid.UUID, domain.Company](redis, "company")
	vacCache := infra_redis.NewRepo[uuid.UUID, domain.Vacancy](redis, "vacancy")
	compListCache := infra_redis.NewRepo[string, list_companies.Response](redis, "companies:list")
	vacListCache := infra_redis.NewRepo[string, list_vacancy.Response](redis, "vacancies:list")
	vacByCompListCache := infra_redis.NewRepo[string, list_vac_by_comp.Response](redis, "vacancies_by_comp:list")

	compGetByIDUc := get_company.NewGetByIDUsecase(compRepo, compCache)
	compListUc := list_companies.NewUsecase(compRepo, compListCache)
	compCreateUc := create_company.NewUsecase(compRepo, memRepo, txManager)
	compAddHrUc := add_member.NewUsecase(memRepo)
	compDeleteMemberUc := delete_member.NewUsecase(memRepo)
	compUpdateMemberUc := update_member.NewUsecase(memRepo)
	compUpdateUc := update_company.NewUsecase(compRepo, memRepo, compCache)
	compDeleteUc := delete_company.NewUsecase(compRepo, memRepo, compCache)

	vacGetByIDUc := get_vacancy.NewUsecase(vacRepo, vacCache)
	vacList := list_vacancy.NewUsecase(vacRepo, vacListCache)
	vacListByComp := list_vac_by_comp.NewUsecase(vacRepo, compRepo, vacByCompListCache)
	vacCreate := create_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager)
	vacUpdate := update_vacancy.NewUsecase(vacRepo, memRepo, vacCache, txManager)
	vacPublish := publish_vacancy.NewUsecase(vacRepo, memRepo)
	vacArchive := archive_vacancy.NewUsecase(vacRepo, memRepo)
	vacDelete := delete_vacancy.NewUsecase(vacRepo, memRepo, vacCache)

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
		vacList,
		vacListByComp,
		vacCreate,
		vacUpdate,
		vacPublish,
		vacArchive,
		vacDelete,
	)

	authMiddleware, err := my_middleware.NewAuthMiddleware(conf)
	if err != nil {
		return nil, err
	}

	routerDeps := &delivery_http.RouterDeps{
		CompanyHandler: companyHandler,
		MemberHandler:  memberHandler,
		VacancyHandler: vacancyHandler,
		AuthMiddleware: authMiddleware,
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
