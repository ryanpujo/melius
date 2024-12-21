package models

type Credential struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type CredentialPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
