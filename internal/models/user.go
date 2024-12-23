package models

type User struct {
	ID         uint       `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Credential Credential `json:"credential"`
}

type UserPayload struct {
	FirstName         string            `json:"first_name" binding:"required"`
	LastName          string            `json:"last_name" binding:"required"`
	CredentialPayload CredentialPayload `json:"credential" binding:"required"`
}
