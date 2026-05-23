package update

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	compRepo      CompanyRepo
	memberRepo    CompMemberRepo
	outboxWriter  outboxWriter
	txManager     common.TxManager
	vacSearchRepo vacSearchRepo
	cache         CacheRepo
}

func NewUsecase(
	repo CompanyRepo,
	memberRepo CompMemberRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	vacSearchRepo vacSearchRepo,
	cache CacheRepo,
) *Usecase {
	return &Usecase{
		compRepo:      repo,
		memberRepo:    memberRepo,
		cache:         cache,
		outboxWriter:  outboxWriter,
		txManager:     txManager,
		vacSearchRepo: vacSearchRepo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *identity.Identity) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := u.authorize(ctx, req.ID, identity); err != nil {
		return err
	}

	nameUpd := false

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		oldName, err := u.compRepo.UpdateAndGetOldName(ctx, req)
		if err != nil {
			return err
		}

		if req.Name != nil && *req.Name != oldName {
			err = u.createCompanyUpdatedEvent(ctx, req.ID, *req.Name)
			if err != nil {
				return err
			}
			nameUpd = true
		}

		u.cache.Del(ctx, req.ID)
		return nil
	})
	if err != nil {
		return err
	}

	if nameUpd {
		err = u.vacSearchRepo.UpdateCompanyName(ctx, req.ID, *req.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// only admin of company can update
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident *identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	memb, err := u.memberRepo.Get(ctx, ident.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}
	if err != nil {
		return err
	}

	if memb.Role != member.CompanyRoleAdmin {
		return member.ErrInsufficientRoleInCompany
	}

	return nil
}

func (u *Usecase) createCompanyUpdatedEvent(ctx context.Context, compID uuid.UUID, newName string) error {
	ev := company.UpdatedEvent{
		EventID:     uuid.New(),
		CompanyID:   compID,
		CompanyName: newName,
		OccurredAt:  time.Now().UTC(),
	}

	return u.outboxWriter.WriteCompanyUpdated(ctx, ev)
}
