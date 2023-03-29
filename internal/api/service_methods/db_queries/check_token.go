package db_queries

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetUserSession(accessToken string) (models.User, error) {
	var user models.User
	query := `
		SELECT u.id,
			   u.company_id,
			   u.local_iso,
			   u.active_flag,
			   json_agg(ru.role_id) role_ids
		FROM sessions s
				 JOIN public.users u ON u.id = s.user_id
				 LEFT JOIN roles_users ru on u.id = ru.user_id
		WHERE access_token = $1
		  AND access_expiration_date > (now() at time zone 'utc')
		GROUP BY u.id,
				 u.company_id,
				 u.local_iso,
				 u.active_flag
	`
	err := postgres.DB.Get(&user, query, accessToken)

	return user, err
}
