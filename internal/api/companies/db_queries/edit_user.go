package dbqueries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

// IsExtIDAlreadyExists returns true if external id already exists in database
func IsExtIDAlreadyExists(tx *sqlx.Tx, extID string) bool {
	query := `SELECT * FROM users WHERE ext_id = $1`
	_, err := tx.Query(query, extID)
	return err == nil
}

func DeleteUserRoles(tx *sqlx.Tx, userID int64) error {
	query := `DELETE FROM roles_users WHERE user_id = $1`
	_, err := tx.Exec(query, userID)
	return err
}

func GetUsersWithPhone(tx *sqlx.Tx, phone string) ([]model.User, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.first_name,
			u.last_name,
			u.middle_name,
			u.position_id,
			u.manager_id,
			u.phone,
			u.active_flag,
			u.company_id,
		FROM
			users
				LEFT JOIN roles_users ru ON ru.user_id = u.id
		WHERE
			phone = $1`
	var users []model.User
	err := tx.Select(&users, query, phone)
	return users, err
}

// GetUserStructures returns user structures
func GetUserStructures(tx *sqlx.Tx, userID int64) ([]models.UserStructure, error) {
	query := `
		SELECT *
		FROM users_company_structures
		WHERE user_id = $1`
	var structures []models.UserStructure
	err := tx.Select(&structures, query, userID)
	return structures, err
}

// GetUserFields returns user fields
func GetUserFields(tx *sqlx.Tx, userID int64) ([]models.UserField, error) {
	query := `
		SELECT *
		FROM users_fields
		WHERE user_id = $1`
	var fields []models.UserField
	err := tx.Select(&fields, query, userID)
	return fields, err
}

func UpdateUser(tx *sqlx.Tx, user *model.User) error {
	query := `
		UPDATE users
		SET	
			first_name = :first_name,
			last_name = :last_name,
			middle_name = :middle_name,
			local_iso = :local_iso,
			phone = :phone,
			image_id = :image_id,
			ext_id = :ext_id,
			position_id = :position_id,
			manager_id = :manager_id,
			active_flag = :active_flag,
			company_id = :company_id
		WHERE id = :id`
	_, err := tx.NamedExec(query, user)
	return err
}
