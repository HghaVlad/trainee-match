package update_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update/mocks"
)

type FakeTxManager struct {
	called bool
}

func (m *FakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	m.called = true
	return fn(ctx)
}

type testDeps struct {
	memRepo   *mocks.MockCompMemberRepo
	vacRepo   *mocks.MockVacancyRepo
	outbox    *mocks.MockoutboxWriter
	vacCache  *mocks.MockCacheRepo
	txManager *FakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		memRepo:   mocks.NewMockCompMemberRepo(ctrl),
		vacRepo:   mocks.NewMockVacancyRepo(ctrl),
		outbox:    mocks.NewMockoutboxWriter(ctrl),
		vacCache:  mocks.NewMockCacheRepo(ctrl),
		txManager: new(FakeTxManager),
	}
}

func NewUC(deps *testDeps) *update.Usecase {
	return update.NewUsecase(deps.vacRepo, deps.memRepo, deps.outbox, deps.vacCache, deps.txManager)
}

type vacMatcher struct {
	expectedVac *vacancy.Vacancy
}

func (m vacMatcher) Matches(x any) bool {
	v, ok := x.(*vacancy.Vacancy)
	if !ok {
		return false
	}

	return m.expectedVac.ID == v.ID && m.expectedVac.Title == v.Title &&
		m.expectedVac.CompanyID == v.CompanyID &&
		m.expectedVac.WorkFormat == v.WorkFormat &&
		m.expectedVac.SalaryTo == v.SalaryTo
}

func (m vacMatcher) String() string {
	return "match vacancy"
}

type vacUpdatedEvMatcher struct {
	expectedEv vacancy.UpdatedEvent
}

func (m vacUpdatedEvMatcher) Matches(x any) bool {
	ev, ok := x.(vacancy.UpdatedEvent)
	if !ok {
		return false
	}

	return m.expectedEv.VacancyID == ev.VacancyID && m.expectedEv.Title == ev.Title
}

func (m vacUpdatedEvMatcher) String() string {
	return "match vacancy updated event"
}

func TestUsecase_Execute_OK_CreatesUpdatedEvent(t *testing.T) {
	vID := uuid.New()
	cID := uuid.New()
	newTitle := "New Go dev"
	newDesc := "New Desc"
	sameSalary := 10000

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	req := &update.Request{
		VacancyID:   vID,
		CompanyID:   cID,
		Title:       &newTitle,
		Description: &newDesc,
	}

	vac := &vacancy.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          "Old title",
		Description:    "Old desc",
		SalaryTo:       &sameSalary,
		IsPaid:         true,
		Status:         vacancy.StatusPublished,
		WorkFormat:     vacancy.WorkFormatHybrid,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}

	expectedVac := &vacancy.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          newTitle,
		Description:    newDesc,
		SalaryTo:       &sameSalary,
		IsPaid:         true,
		Status:         vacancy.StatusPublished,
		WorkFormat:     vacancy.WorkFormatHybrid,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}

	expectedEv := vacancy.UpdatedEvent{
		VacancyID: vID,
		Title:     newTitle,
	}

	deps := setup(t)

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, cID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil)

	deps.vacRepo.EXPECT().GetByIDForUpdate(gomock.Any(), vID, cID).Return(vac, nil)

	deps.vacRepo.EXPECT().Update(gomock.Any(), vacMatcher{expectedVac}).
		Return(nil)

	deps.outbox.EXPECT().WriteVacancyUpdated(gomock.Any(), vacUpdatedEvMatcher{expectedEv}).
		Return(nil)

	deps.vacCache.EXPECT().Del(gomock.Any(), vID)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), req, ident)

	require.NoError(t, err)
	assert.True(t, deps.txManager.called)
}

func TestUsecase_Execute_OK_NoEvent(t *testing.T) {
	vID := uuid.New()
	cID := uuid.New()
	oldSalary := 2000
	newSalary := 10000

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	req := &update.Request{
		VacancyID: vID,
		CompanyID: cID,
		SalaryTo:  &newSalary,
	}

	vac := &vacancy.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          "Old title",
		Description:    "Old desc",
		SalaryTo:       &oldSalary,
		IsPaid:         true,
		Status:         vacancy.StatusPublished,
		WorkFormat:     vacancy.WorkFormatHybrid,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}

	expectedVac := &vacancy.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          "Old title",
		Description:    "Old desc",
		SalaryTo:       &newSalary,
		IsPaid:         true,
		Status:         vacancy.StatusPublished,
		WorkFormat:     vacancy.WorkFormatHybrid,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}

	deps := setup(t)

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, cID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil)

	deps.vacRepo.EXPECT().GetByIDForUpdate(gomock.Any(), vID, cID).Return(vac, nil)

	deps.vacRepo.EXPECT().Update(gomock.Any(), vacMatcher{expectedVac}).
		Return(nil)

	deps.vacCache.EXPECT().Del(gomock.Any(), vID)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), req, ident)

	require.NoError(t, err)
	assert.True(t, deps.txManager.called)
}

func TestUsecase_Execute_ConflictingChangesToExistingState(t *testing.T) {
	vID := uuid.New()
	cID := uuid.New()
	oldSalaryFrom := 10000
	oldSalaryTo := 20000
	newSalaryTo := 5000 // request is valid, but new salaryTo < oldSalaryFrom

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	req := &update.Request{
		VacancyID: vID,
		CompanyID: cID,
		SalaryTo:  &newSalaryTo,
	}

	vac := &vacancy.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          "Old title",
		Description:    "Old desc",
		SalaryFrom:     &oldSalaryFrom,
		SalaryTo:       &oldSalaryTo,
		IsPaid:         true,
		Status:         vacancy.StatusPublished,
		WorkFormat:     vacancy.WorkFormatHybrid,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}

	deps := setup(t)

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, cID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil)

	deps.vacRepo.EXPECT().GetByIDForUpdate(gomock.Any(), vID, cID).Return(vac, nil)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), req, ident)

	require.ErrorIs(t, err, vacancy.ErrInvalidSalaryRange)
}

func TestUsecase_Execute_VacancyNotFound(t *testing.T) {
	vID := uuid.New()
	cID := uuid.New()

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	req := &update.Request{
		VacancyID: vID,
		CompanyID: cID,
		Title:     ptr("new title"),
	}

	deps := setup(t)

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, cID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil)

	deps.vacRepo.EXPECT().GetByIDForUpdate(gomock.Any(), vID, cID).
		Return(nil, vacancy.ErrVacancyNotFound)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), req, ident)

	require.ErrorIs(t, err, vacancy.ErrVacancyNotFound)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	vID := uuid.New()
	cID := uuid.New()

	req := &update.Request{
		VacancyID: vID,
		CompanyID: cID,
		Title:     ptr("new title"),
	}

	t.Run("not hr role", func(t *testing.T) {
		deps := setup(t)
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	})

	t.Run("not company member", func(t *testing.T) {
		deps := setup(t)
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, cID).
			Return(nil, member.ErrCompanyMemberNotFound)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
	})
}

func ptr[T any](v T) *T {
	return &v
}
