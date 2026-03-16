package get_published_vacancy

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Usecase returns read model of published vacancy (for candidates)
type Usecase struct {
	repo  Repository
	cache CacheRepo
}

func NewUsecase(repo Repository, cache CacheRepo) *Usecase {
	return &Usecase{
		repo:  repo,
		cache: cache,
	}
}

// Execute returns read model of published vacancy (for candidates)
func (u *Usecase) Execute(ctx context.Context, vacancyID uuid.UUID) (*Response, error) {
	vacancy := u.cache.Get(ctx, vacancyID)
	if vacancy != nil {
		return vacancy, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	vacancy, err := u.repo.GetPublishedByID(ctx, vacancyID)
	if err != nil {
		return nil, err
	}

	u.cache.Put(ctx, vacancyID, vacancy, time.Second*300)
	return vacancy, nil
}
