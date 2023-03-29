package db_queries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetUserF(email string) (models.DBUser, error) {
	var user models.DBUser
	query := `
		SELECT id, email, first_name, middle_name, sign_up_token, local_iso
		FROM users u
		WHERE email = $1
	`
	err := postgres.DB.Get(&user, query, email)
	return user, err
}

func UpdateUsersignUpToken(tx *sqlx.Tx, userId int64, signUpToken string) error {

	_, err := tx.Exec(`UPDATE users SET sign_up_token = $1 where id= $2`, signUpToken, userId)
	if err != nil {
		return err
	}

	return nil
}
