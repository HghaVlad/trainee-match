package remove

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

func (u *Usecase) Execute(ctx context.Context, companyID, userID uuid.UUID, identity *identity.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, companyID, identity); err != nil {
			return err
		}

		if identity.UserID == userID {
			if err := u.checkIfCanRemoveYourself(ctx, companyID); err != nil {
				return err
			}
		}

		if err := u.memberRepo.Delete(ctx, userID, companyID); err != nil {
			return err
		}

		return u.createCompMemRemovedEvent(ctx, userID, companyID)
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

func (u *Usecase) createCompMemRemovedEvent(ctx context.Context, userID, companyID uuid.UUID) error {
	ev := member.RemovedEvent{
		EventID:    uuid.New(),
		UserID:     userID,
		CompanyID:  companyID,
		OccurredAt: time.Now().UTC(),
	}

	return u.outboxWriter.WriteCompanyMemberRemoved(ctx, ev)
}

// admin can remove themselves if they are not the only admin in company
func (u *Usecase) checkIfCanRemoveYourself(ctx context.Context, companyID uuid.UUID) error {
	adminCnt, err := u.memberRepo.GetCompanyRoleCount(ctx, companyID, member.CompanyRoleAdmin)
	if err != nil {
		return err
	}

	if adminCnt == 1 {
		return member.ErrCantRemoveYourself
	}

	return nil
}
