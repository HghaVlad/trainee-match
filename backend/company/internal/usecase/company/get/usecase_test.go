package get_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
)

type CompanyRepoMock struct {
	mock.Mock
}

func (m *CompanyRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	args := m.Called(ctx, id)

	if c := args.Get(0); c != nil {
		return c.(*company.Company), args.Error(1)
	}

	return nil, args.Error(1)
}

type CacheRepoMock struct {
	mock.Mock
}

func (m *CacheRepoMock) Get(ctx context.Context, id uuid.UUID) *company.Company {
	args := m.Called(ctx, id)

	if c := args.Get(0); c != nil {
		return c.(*company.Company)
	}

	return nil
}

func (m *CacheRepoMock) Put(
	ctx context.Context,
	id uuid.UUID,
	company *company.Company,
	ttl time.Duration,
) {
	m.Called(ctx, id, company, ttl)
}

// Table tests
func TestGetByIdUsecase_Execute(t *testing.T) {
	id := uuid.New()

	now := time.Now()
	comp := &company.Company{
		ID:               id,
		Name:             "Acme",
		OpenVacanciesCnt: 3,
		Description:      ptr("desc"),
		Website:          ptr("site"),
		LogoKey:          ptr("logo.png"),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name     string
		setup    func(repo *CompanyRepoMock, cache *CacheRepoMock)
		wantErr  bool
		wantName string
	}{
		{
			name: "cache hit",
			setup: func(repo *CompanyRepoMock, cache *CacheRepoMock) {
				cache.
					On("Get", mock.Anything, id).
					Return(comp)

				repo.
					On("GetByID", mock.Anything, id).
					Maybe()
			},
			wantName: "Acme",
		},
		{
			name: "cache miss -> repo hit",
			setup: func(repo *CompanyRepoMock, cache *CacheRepoMock) {
				cache.
					On("Get", mock.Anything, id).
					Return(nil)

				repo.
					On("GetByID", mock.Anything, id).
					Return(comp, nil)

				cache.
					On("Put", mock.Anything, id, comp, 5*time.Minute).
					Once()
			},
			wantName: "Acme",
		},
		{
			name: "repo error",
			setup: func(repo *CompanyRepoMock, cache *CacheRepoMock) {
				cache.
					On("Get", mock.Anything, id).
					Return(nil)

				repo.
					On("GetByID", mock.Anything, id).
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(CompanyRepoMock)
			cache := new(CacheRepoMock)

			tt.setup(repo, cache)

			uc := get.NewGetByIDUsecase(repo, cache)

			resp, err := uc.Execute(context.Background(), id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantName, resp.Name)

			repo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
