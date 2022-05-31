package models

type SessionPayload struct {
	Token string `json:"token" validate:"required"`
}
