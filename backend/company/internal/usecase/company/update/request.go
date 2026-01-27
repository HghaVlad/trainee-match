package update_company

import "github.com/google/uuid"

type Request struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Website     *string
}
