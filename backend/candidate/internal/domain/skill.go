package domain

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrSkillNotFound = errors.New("skill not found")
)

type Skill struct {
	ID   uuid.UUID
	Name string
}
