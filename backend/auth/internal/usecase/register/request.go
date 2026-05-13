package register

import (
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

type Request struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name"  validate:"required,min=2,max=50"`
	Email     string `json:"email"      validate:"required,email,max=254"`
	Username  string `json:"username"   validate:"required,min=3,max=50,alphanum"`
	Password  string `json:"password"   validate:"required,min=8"`
	Role      string `json:"role"       validate:"required,oneof=Candidate Company"`
}

func (r *Request) Validate(validate *validator.Validate) error {
	if validate == nil {
		return errors.New("validator is nil")
	}
	if err := validate.Struct(r); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			for _, fe := range verrs {
				switch fe.Field() {
				case "Email":
					return domain.ErrInvalidEmail
				case "Password":
					return domain.ErrInvalidPassword
				case "FirstName", "LastName":
					return domain.ErrInvalidName
				}
			}
		}
		return err
	}
	return nil
}
