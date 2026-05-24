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
	cache         CacheRepo
}

func NewUsecase(
	repo CompanyRepo,
	memberRepo CompMemberRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	cache CacheRepo,
) *Usecase {
	return &Usecase{
		compRepo:      repo,
		memberRepo:    memberRepo,
		cache:         cache,
		outboxWriter:  outboxWriter,
		txManager:     txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *identity.Identity) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	compID := req.ID

	if err := u.authorize(ctx, compID, identity); err != nil {
		return err
	}

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		oldName, err := u.compRepo.UpdateAndGetOldName(ctx, req)
		if err != nil {
			return err
		}

		return u.createCompanyUpdatedEvent(ctx, compID, oldName, req)
	})
	if err != nil {
		return err
	}

	u.cache.Del(ctx, compID)
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

func (u *Usecase) createCompanyUpdatedEvent(ctx context.Context, compID uuid.UUID, oldName string, req *Request) error {
	if req.Name == nil || *req.Name == oldName {
		return nil
	}

	ev := company.UpdatedEvent{
		EventID:     uuid.New(),
		CompanyID:   compID,
		CompanyName: *req.Name,
		OccurredAt:  time.Now().UTC(),
	}

	return u.outboxWriter.WriteCompanyUpdated(ctx, ev)
}
