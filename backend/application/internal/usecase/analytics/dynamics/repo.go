package dynamics

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type historyRepo interface {
	GetDynamicsBucketsByCompany(
		ctx context.Context,
		compID uuid.UUID,
		from, to time.Time,
		interval Interval,
	) ([]Bucket, error)

	GetDynamicsBucketsByVacancy(
		ctx context.Context,
		vacID uuid.UUID,
		from, to time.Time,
		interval Interval,
	) ([]Bucket, error)
}

type companyMemberRepo interface {
	IsMember(
		ctx context.Context,
		userID uuid.UUID,
		companyID uuid.UUID,
	) (bool, error)
}

type vacancyProjRepo interface {
	CheckHrAccess(
		ctx context.Context,
		userID, vacancyID uuid.UUID,
	) (uuid.UUID, error)
}
