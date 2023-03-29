package models

type ResetPassword struct {
	Email string `json:"email" binding:"email,max=255"`
}
