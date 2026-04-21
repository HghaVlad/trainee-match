package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	grpc2 "github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/grpc"

	candidatev1 "github.com/HghaVlad/trainee-match/backend/contracts/go/candidate/v1"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	myhttp "github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate_by_user_id"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_skill"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume"
)

type App struct {
	httpServer   *http.Server
	grpcServer   *grpc.Server
	grpcListener net.Listener
	Db           *pgxpool.Pool
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

	createCandidateUC := create_candidate.New(candidateRepo)
	updateCandidateUC := update_candidate.New(candidateRepo)
	getCandidateByUserIdUC := get_candidate_by_user_id.New(candidateRepo)

	getResumeUC := get_resume.New(resumeRepo, candidateRepo)
	createResumeUC := create_resume.New(resumeRepo, skillRepo, candidateRepo)
	updateResumeUC := update_resume.New(resumeRepo, skillRepo, candidateRepo)

	getSkillUC := get_skill.New(skillRepo)

	candidateHandler := handlers.NewCandidate(createCandidateUC, updateCandidateUC, getCandidateByUserIdUC)
	resumeHandler := handlers.NewResume(createResumeUC, getResumeUC, updateResumeUC)
	skillHandler := handlers.NewSkill(getSkillUC)
	authMiddleware := auth.NewMiddleware(conf.JWKUrl)

	router := myhttp.NewRouter(myhttp.NewRouterDeps(authMiddleware, candidateHandler, resumeHandler, skillHandler))

	httpServer := &http.Server{
		Addr:         conf.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	grpcServer := grpc.NewServer()
	candidateService := grpc2.NewCandidateService(getCandidateByUserIdUC, getResumeUC)
	candidatev1.RegisterCandidateServiceServer(grpcServer, candidateService)
	grpcLis, err := net.Listen("tcp", conf.GrpcAddr)
	if err != nil {
		return nil, err
	}

	return &App{
		httpServer:   httpServer,
		Db:           pgPool,
		grpcServer:   grpcServer,
		grpcListener: grpcLis,
	}, nil
}

func (app *App) Run() error {
	errCh := make(chan error, 2)

	go func() {
		err := app.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	go func() {
		if err := app.grpcServer.Serve(app.grpcListener); err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	for {
		err := <-errCh
		if err != nil {
			return err
		}
	}
}

func (app *App) Shutdown(ctx context.Context) {
	if app.httpServer != nil {
		if err := app.httpServer.Shutdown(ctx); err != nil {
			slog.Error("http shutdown error", "error", err)
		}
	}

	if app.grpcServer != nil {
		done := make(chan struct{})
		go func() {
			app.grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
		case <-ctx.Done():
			app.grpcServer.Stop()
		}
	}

	if app.grpcListener != nil {
		if err := app.grpcListener.Close(); err != nil {
			slog.Error("grpc listener close error", "error", err)
		}
	}

	if app.Db != nil {
		app.Db.Close()
	}
}
