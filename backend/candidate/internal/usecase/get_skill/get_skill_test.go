package get_skill

import (
	"context"
	"errors"
	"testing"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_skill/mocks"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

var (
	ErrDb = errors.New("db error")
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	domainSkill := domain.Skill{ID: id, Name: "Go"}

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.SkillRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(repo *mocks.SkillRepo) {
				repo.On("GetByID", ctx, id).Return(domainSkill, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(repo *mocks.SkillRepo) {
				repo.On("GetByID", ctx, id).Return(domain.Skill{}, domain.ErrSkillNotFound).Once()
			},
			expectedError: domain.ErrSkillNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(repo *mocks.SkillRepo) {
				repo.On("GetByID", ctx, id).Return(domain.Skill{}, ErrDb).Once()
			},
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.SkillRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := New(repo)
			resp, err := uc.Execute(ctx, GetByIdRequest{ID: id})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected %v got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, domainSkill.ID, resp.ID)
				require.Equal(t, domainSkill.Name, resp.Name)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestExecuteList(t *testing.T) {
	ctx := context.Background()
	domainSkills := []domain.Skill{{ID: uuid.New(), Name: "Go"}, {ID: uuid.New(), Name: "Python"}}

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.SkillRepo)
		expectedError error
	}{
		{
			name: "valid list",
			mockSetup: func(repo *mocks.SkillRepo) {
				repo.On("List", ctx).Return(domainSkills, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "repo error",
			mockSetup: func(repo *mocks.SkillRepo) {
				repo.On("List", ctx).Return(nil, ErrDb).Once()
			},
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.SkillRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := New(repo)
			resp, err := uc.ExecuteList(ctx, ListRequest{})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected %v got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(domainSkills), len(resp))
			}

			repo.AssertExpectations(t)
		})
	}
}
