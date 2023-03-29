package db_queries

import (
	"time"

	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func SetChangePasswordToken(changePassToken string, tokenExpirationDate time.Time, userId int64) error {
	_, err := postgres.DB.Exec(`
		UPDATE users
		SET change_password_token = $1, change_password_token_expiration_date = $2
		WHERE id = $3
	`, changePassToken, tokenExpirationDate, userId)
	return err
}
