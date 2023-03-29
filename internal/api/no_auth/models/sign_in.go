package models

import (
	"database/sql"

	"gitlab.assistagro.com/back/back.auth.go/internal/models"
)

type SignINRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type DBUser struct {
	models.User

	PasswordHash string         `db:"password_hash"`
	SignUPToken  sql.NullString `db:"sign_up_token"`
}
