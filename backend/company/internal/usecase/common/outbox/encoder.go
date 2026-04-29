package outbox

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type Encoder interface {
	VacancyPublishedToBytes(ev vacancy.PublishedEvent) ([]byte, error)
	VacancyArchivedToBytes(ev vacancy.ArchivedEvent) ([]byte, error)
	VacancyUpdatedToBytes(ev vacancy.UpdatedEvent) ([]byte, error)

	CompanyMemberAddedToBytes(ev member.RecruiterAddedEvent) ([]byte, error)
	CompanyMemberRemovedToBytes(ev member.RecruiterRemovedEvent) ([]byte, error)

	CompanyDeletedToBytes(ev company.DeletedEvent) ([]byte, error)
	CompanyUpdatedToBytes(ev company.UpdatedEvent) ([]byte, error)
}
