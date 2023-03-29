package dbqueries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
)

func GetUserImage(tx *sqlx.Tx, userID int64) (models.DBUserImage, error) {
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
		json_agg(r.role_id) AS role_ids,
		i.image_path
	FROM
		users AS u
		LEFT JOIN roles_users AS r ON r.user_id = u.id
		LEFT JOIN images AS i ON u.image_id=i.id
	WHERE u.id = $1
	GROUP BY u.id, u.image_id, i.image_path
	`
	var res models.DBUserImage
	err := tx.Get(&res, query, userID)
	return res, err
}
