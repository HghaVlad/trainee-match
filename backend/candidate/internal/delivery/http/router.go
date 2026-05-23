package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/handlers"
)

type RouterDeps struct {
	authMiddleware   *auth.Middleware
	candidateHandler *handlers.Candidate
	resumeHandler    *handlers.Resume
	skillHandler     *handlers.Skill
	adminHandler     *handlers.Admin
}

func NewRouterDeps(
	authMiddleware *auth.Middleware,
	candidateHandler *handlers.Candidate,
	resumeHandler *handlers.Resume,
	skillHandler *handlers.Skill,
	adminHandler *handlers.Admin,
) *RouterDeps {
	return &RouterDeps{
		authMiddleware:   authMiddleware,
		candidateHandler: candidateHandler,
		resumeHandler:    resumeHandler,
		skillHandler:     skillHandler,
		adminHandler:     adminHandler,
	}
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Route("/api/v1/candidate", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(deps.authMiddleware.Handler)
			r.Get("/me", deps.candidateHandler.GetMe)
			r.Post("/", deps.candidateHandler.CreateCandidate)
			r.Patch("/", deps.candidateHandler.UpdateCandidate)
		})
	})

	// Resume routes
	router.Route("/api/v1/resume", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(deps.authMiddleware.Handler)
			r.Post("/", deps.resumeHandler.CreateResume)
			r.Get("/", deps.resumeHandler.ListResumes)
			r.Get("/{id}", deps.resumeHandler.GetResume)
			r.Patch("/{id}", deps.resumeHandler.UpdateResume) // Changed from PUT to PATCH
			r.Delete("/{id}", deps.resumeHandler.DeleteResume)
		})
	})

	// Skill routes
	router.Route("/api/v1/skill", func(r chi.Router) {
		r.Get("/{id}", deps.skillHandler.GetSkill)
		r.Get("/list", deps.skillHandler.ListSkills)
		// Note: Create skill functionality has been removed as per requirements
	})

	// Admin routes
	router.Route("/api/v1/admin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(deps.authMiddleware.Handler)
			r.Use(auth.IsAdmin)
			r.Get("/candidates", deps.adminHandler.GetCandidates)
			r.Get("/candidates/{id}", deps.adminHandler.GetCandidate)
			r.Get("/candidates/{id}/resumes", deps.adminHandler.GetCandidateResumes)
			r.Get("/resumes/{id}", deps.adminHandler.GetResume)
			r.Post("/resumes/{id}/archive", deps.adminHandler.ArchiveResume)
			r.Post("/skills", deps.adminHandler.AddSkill)
			r.Delete("/skills/{id}", deps.adminHandler.DeleteSkill)
		})
	})

	router.Get("/swagger/*", handlers.SwaggerHandler)

	return router
}
