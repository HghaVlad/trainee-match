package archive_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	vacRepo    VacancyRepo
	memberRepo CompMemberRepo
}

func NewUsecase(vacRepo VacancyRepo, memberRepo CompMemberRepo) *Usecase {
	return &Usecase{vacRepo: vacRepo, memberRepo: memberRepo}
}

func (u *Usecase) Execute(
	ctx context.Context,
	compID uuid.UUID,
	vacID uuid.UUID,
	identity uc_common.Identity,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO: think about cache here
	if err := u.authorize(ctx, compID, identity); err != nil {
		return err
	}

	err := u.vacRepo.Archive(ctx, compID, vacID)
	return err
}

// only member of company can archive vacancy
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, identity uc_common.Identity) error {
	if identity.Role != uc_common.RoleHR {
		return domain_errors.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, identity.UserID, companyID)
	if errors.Is(err, domain_errors.ErrCompanyMemberNotFound) {
		return domain_errors.ErrCompanyMemberRequired
	}

	return err
}
