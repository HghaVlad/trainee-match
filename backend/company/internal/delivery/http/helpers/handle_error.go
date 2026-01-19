package helpers

import (
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
)

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain_errors.ErrCompanyNotFound):
		RespondError(w, http.StatusNotFound, err)

	default:
		RespondError(w, http.StatusInternalServerError, err)
	}
}
