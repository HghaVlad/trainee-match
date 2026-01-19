package handlers

import (
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/get_profile"
)

type ProfileHandler struct {
	getByID *get_profile.GetByIDUsecase
}

func NewProfileHandler(
// getByID *get_profile.GetByIDUsecase,

) *ProfileHandler {

	return &ProfileHandler{
		//getByID: getByID,
	}
}

func (g *ProfileHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
