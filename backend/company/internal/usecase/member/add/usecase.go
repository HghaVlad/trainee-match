package add

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	memberRepo   CompanyMemberRepo
	outboxWriter outboxWriter
	txManager    common.TxManager
}

func NewUsecase(
	memberRepo CompanyMemberRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		memberRepo:   memberRepo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *identity.Identity) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, req.CompanyID, identity); err != nil {
			return err
		}

		memb := &member.CompanyMember{
			UserID:    req.UserID,
			CompanyID: req.CompanyID,
			Role:      req.Role,
		}

		if err := u.memberRepo.Create(ctx, memb); err != nil {
			return err
		}

		return u.createCompMemAddedEvent(ctx, memb)
	})
}

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

func (u *Usecase) createCompMemAddedEvent(ctx context.Context, mem *member.CompanyMember) error {
	ev := member.AddedEvent{
		EventID:    uuid.New(),
		UserID:     mem.UserID,
		CompanyID:  mem.CompanyID,
		Role:       mem.Role,
		OccurredAt: time.Now().UTC(),
	}

	return u.outboxWriter.WriteCompanyMemberAdded(ctx, ev)
}
