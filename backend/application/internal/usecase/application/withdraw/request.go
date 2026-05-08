package withdraw

import "github.com/google/uuid"

type Request struct {
	AppID   uuid.UUID
	Comment *string
}
