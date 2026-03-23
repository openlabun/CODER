package postlogin

type RequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
