package utilities

type RegistrationResponse struct {
	ID      uint   `json:"id"`
	Token   string `json:"token"`
	Err     error  `json:"err"`
	Message string `json:"message"`
}
