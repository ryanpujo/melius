package models

type Credential struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type CredentialPayload struct {
	Email    string `json:"email" binding:"email,required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
