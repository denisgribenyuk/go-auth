package models

type ResendEmail struct {
	Email string `json:"email" binding:"required,email"`
}
