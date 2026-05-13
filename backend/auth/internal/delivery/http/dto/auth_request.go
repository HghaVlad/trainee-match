package dto

// RegisterUserRequest registration request
// @name RegisterUserRequest
type RegisterUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name"  validate:"required,min=2,max=50"`
	Email     string `json:"email"      validate:"required,email,max=254"`
	Username  string `json:"username"   validate:"required,min=3,max=50,alphanum"`
	Password  string `json:"password"   validate:"required,min=8"`
	Role      string `json:"role"       validate:"required,oneof=Candidate Company"`
}

// LoginRequest login request
// @name LoginRequest
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest refresh token request
// @name RefreshTokenRequest
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
