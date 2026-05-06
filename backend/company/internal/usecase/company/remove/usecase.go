package remove

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
	compRepo     CompanyRepo
	memberRepo   CompMemberRepo
	outboxWriter outboxWriter
	txManager    common.TxManager
	cache        CacheRepo
}

func NewUsecase(
	repo CompanyRepo,
	memberRepo CompMemberRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	cache CacheRepo,
) *Usecase {
	return &Usecase{
		compRepo:     repo,
		memberRepo:   memberRepo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
		cache:        cache,
	}
}

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID, identity *identity.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		err := u.authorize(ctx, id, identity)
		if err != nil {
			return err
		}

		err = u.compRepo.Delete(ctx, id)
		if err != nil {
			return err
		}

		err = u.createCompanyDeletedEvent(ctx, id)
		if err != nil {
			return err
		}

		u.cache.Del(ctx, id)
		return nil
	})
}

// only admin of company can delete or admin of the platform
func (u *Usecase) authorize(ctx context.Context, id uuid.UUID, ident *identity.Identity) error {
	if ident.Role == identity.RoleAdmin {
		return nil
	}

	if ident.Role != identity.RoleHR {
		return identity.ErrInsufficientRole
	}

	memb, err := u.memberRepo.Get(ctx, ident.UserID, id)
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

func (u *Usecase) createCompanyDeletedEvent(ctx context.Context, compID uuid.UUID) error {
	ev := company.DeletedEvent{
		EventID:    uuid.New(),
		CompanyID:  compID,
		OccurredAt: time.Now(),
	}

	return u.outboxWriter.WriteCompanyDeleted(ctx, ev)
}
