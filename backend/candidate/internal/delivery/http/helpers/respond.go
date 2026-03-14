package helpers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, dto.ErrorResponse{
		Message: message,
	})
}

// RespondErrorSmart responds with error and select status itself
func RespondErrorSmart(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, domain.ErrCandidateNotFound),
		errors.Is(err, domain.ErrResumeNotFound),
		errors.Is(err, domain.ErrSkillNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrCandidateAlreadyExists),
		errors.Is(err, domain.ErrTelegramAlreadyExists),
		errors.Is(err, domain.ErrPhoneAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, domain.ErrInvalidPhoneFormat),
		errors.Is(err, domain.ErrInvalidTelegramFormat),
		errors.Is(err, domain.ErrInvalidCityFormat),
		errors.Is(err, domain.ErrBirthdayInFuture),
		errors.Is(err, domain.ErrInvalidResumeName),
		errors.Is(err, domain.ErrInvalidResumeStatus),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidEmailFormat),
		errors.Is(err, domain.ErrDateOfBirthInFuture),
		errors.Is(err, domain.ErrInvalidCitizenship),
		errors.Is(err, domain.ErrInvalidEducationEntry),
		errors.Is(err, domain.ErrInvalidWorkExperienceEntry),
		errors.Is(err, domain.ErrInvalidPortfolioLink),
		errors.Is(err, domain.ErrInvalidSkillName):
		status = http.StatusBadRequest
	}

	RespondError(w, status, err.Error())
}
