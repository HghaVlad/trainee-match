package app

import (
	"context"
	"errors"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	myhttp "github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate_by_user_id"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_skill"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server *http.Server
	Db     *pgxpool.Pool
}

func Build(conf *config.Config) (*App, error) {
	pgPool, err := postgres.Connect(context.Background(), &conf.Db)
	if err != nil {
		return nil, err
	}
	err = postgres.Migrate(conf.Db.GetPostgresURL())
	if err != nil {
		return nil, err
	}

	candidateRepo := repository.NewCandidateRepo(pgPool)
	resumeRepo := repository.NewResumeRepo(pgPool)
	skillRepo := repository.NewSkillRepo(pgPool)

	getCandidateUC := get_candidate.New(candidateRepo)
	createCandidateUC := create_candidate.New(candidateRepo)
	updateCandidateUC := update_candidate.New(candidateRepo)
	getCandidateByUserIdUC := get_candidate_by_user_id.New(candidateRepo)

	getResumeUC := get_resume.New(resumeRepo)
	createResumeUC := create_resume.New(resumeRepo, skillRepo)
	updateResumeUC := update_resume.New(resumeRepo, skillRepo)

	getSkillUC := get_skill.New(skillRepo)

	candidateHandler := handlers.NewCandidate(getCandidateUC, createCandidateUC, updateCandidateUC, getCandidateByUserIdUC)
	resumeHandler := handlers.NewResume(createResumeUC, getResumeUC, updateResumeUC, getCandidateByUserIdUC)
	skillHandler := handlers.NewSkill(getSkillUC)
	authMiddleware := auth.NewMiddleware(conf.JWKUrl)

	router := myhttp.NewRouter(myhttp.NewRouterDeps(authMiddleware, candidateHandler, resumeHandler, skillHandler))

	httpServer := &http.Server{
		Addr:         conf.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		server: httpServer,
		Db:     pgPool,
	}, nil
}

func (app *App) Run() error {

	slog.Info("Server started")
	err := app.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http listening server err: %w \n", err)
	}
	return err
}

func (app *App) Shutdown(ctx context.Context) {
	err := app.server.Shutdown(ctx)
	if err != nil {
		slog.Error("shutdown error", err)
	}
	slog.Info("Server stopped")
	app.Db.Close()
}
