package models

type User struct {
	ID         uint       `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Credential Credential `json:"credential"`
}

type UserPayload struct {
	ID                uint              `json:"id"`
	FirstName         string            `json:"first_name"`
	LastName          string            `json:"last_name"`
	CredentialPayload CredentialPayload `json:"credential"`
}
