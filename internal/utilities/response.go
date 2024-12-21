package utilities

type RegistrationResponse struct {
	ID      uint   `json:"id"`
	Err     error  `json:"-"`
	Message string `json:"message"`
}
