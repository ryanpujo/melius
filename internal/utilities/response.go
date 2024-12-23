package utilities

type RegistrationResponse struct {
	ID      uint   `json:"id"`
	Token   string `json:"token"`
	Err     string `json:"err"`
	Message string `json:"message"`
}
