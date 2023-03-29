package db_queries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/internal/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetSession(tx *sqlx.Tx, refreshToken string) (*models.Session, error) {
	var session models.Session
	query := `
		SELECT *
		FROM sessions
		WHERE refresh_token = $1
		  AND refresh_expiration_date > (now() at time zone 'utc')
		  AND user_id IS NOT NULL
	`
	err := tx.Get(&session, query, refreshToken)
	return &session, err
}

func GetUserByID(id int64) (*models.User, error) {
	var user models.User
	query := `
		SELECT u.id,
			   u.email,
		       u.company_id,
		       u.active_flag
		FROM users u
		WHERE id = $1
	`
	err := postgres.DB.Get(&user, query, id)
	return &user, err
}

func UpdateSession(tx *sqlx.Tx, session *models.Session) error {
	query := `
		UPDATE sessions
		SET access_token=:access_token,
			access_expiration_date=:access_expiration_date,
			refresh_token=:refresh_token,
			refresh_expiration_date=:refresh_expiration_date
		WHERE id = :id
	`
	_, err := tx.NamedExec(query, session)
	return err
}
