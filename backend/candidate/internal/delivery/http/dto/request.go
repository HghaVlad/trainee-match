package dto

type CandidateCreateRequest struct {
	Phone    string `json:"phone"`
	Telegram string `json:"telegram"`
	City     string `json:"city"`
	Birthday string `json:"birthday"`
}
