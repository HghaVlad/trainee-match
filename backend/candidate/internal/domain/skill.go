package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrSkillNotFound    = errors.New("skill not found")
	ErrInvalidSkillName = errors.New("invalid skill name")
)

type Skill struct {
	ID   uuid.UUID
	Name string
}

func (s Skill) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return ErrInvalidSkillName
	}
	return nil
}
