package dbqueries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetUserByID(id int64) (*model.User, error) {
	query := `
	SELECT
        u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.middle_name,
		u.local_iso,
		u.company_id,
		u.phone,
		u.image_id,
		u.ext_id,
		u.manager_id,
		u.position_id,
		json_agg(r.role_id) AS role_ids
	FROM
			users u
			LEFT JOIN roles_users r ON r.user_id = u.id
	WHERE
		u.id=$1        
	GROUP BY u.id;
	`
	var user model.User
	err := postgres.DB.Get(&user, query, id)

	return &user, err
}

func GetUserByToken(accessToken string) (*model.User, error) {
	query := `
	SELECT
        u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.middle_name,
		u.local_iso,
		u.company_id,
		u.phone,
		u.image_id,
		u.ext_id,
		u.manager_id,
		u.position_id,
		json_agg(r.role_id) AS role_ids
	FROM
			users u
			LEFT JOIN roles_users r ON r.user_id = u.id,
			sessions s
	WHERE
		u.id = s.user_id
		AND access_token=$1        
	GROUP BY u.id;
	`

	var user model.User
	err := postgres.DB.Get(&user, query, accessToken)
	return &user, err
}

func GetUserWithPhone(tx *sqlx.Tx, phone string) ([]model.User, error) {
	query := `
	SELECT
        u.id,
		u.email,
		u.first_name,
		u.last_name,
		u.middle_name,
		u.local_iso,
		u.company_id,
		u.phone,
		u.image_id,
		u.ext_id,
		u.manager_id,
		u.position_id,
		json_agg(r.role_id) AS role_ids
	FROM
			users u
			LEFT JOIN roles_users r ON r.user_id = u.id
	WHERE
		phone=$1        
	GROUP BY u.id;
	`

	var users []model.User
	err := postgres.DB.Select(&users, query, phone)
	return users, err
}
