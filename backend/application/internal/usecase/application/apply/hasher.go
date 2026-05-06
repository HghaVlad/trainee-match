package apply

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type snapIDGetter interface {
	GetDeterministicAppSnapshotID(resumeData projection.ResumeData, candidateProj projection.Candidate) (uuid.UUID, error)
}
