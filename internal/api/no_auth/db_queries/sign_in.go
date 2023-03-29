package db_queries

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	internal_models "gitlab.assistagro.com/back/back.auth.go/internal/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetUser(email string) (models.DBUser, error) {
	var user models.DBUser
	query := `
		SELECT u.id,
			   u.email,
			   u.password_hash,
		       u.local_iso,
			   u.sign_up_token,
		       u.company_id,
		       u.active_flag
		FROM users u
		WHERE email = $1
	`
	err := postgres.DB.Get(&user, query, email)
	return user, err
}

func CreateSession(session internal_models.Session) error {
	query := `
		INSERT INTO sessions (user_id, access_token, refresh_token, access_expiration_date, refresh_expiration_date)
		VALUES (:user_id, :access_token, :refresh_token, :access_expiration_date, :refresh_expiration_date)
	`
	_, err := postgres.DB.NamedExec(query, session)
	return err
}
