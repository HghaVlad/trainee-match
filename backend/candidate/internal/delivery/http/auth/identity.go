package auth

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Username  string
}
