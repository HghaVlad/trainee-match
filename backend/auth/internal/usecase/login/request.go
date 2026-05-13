package login

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
