package dbqueries

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
)

func CheckUserEmailPhoneExtIDExists(tx *sqlx.Tx, email string, phone string, extID *string) (bool, error) {

	var res []int64

	query := `
		SELECT
			id
		FROM
			users
		WHERE
			email = $1
			OR phone = $2
	`

	if extID != nil {
		query += " OR ext_id = $3"
		err := tx.Select(&res, query, email, phone, *extID)
		return len(res) > 0, err
	}

	err := tx.Select(&res, query, email, phone)

	return len(res) > 0, err
}

// Insert user inserts new user into DB and returns its ID
func InsertUser(tx *sqlx.Tx, userObj *models.PostNewUserReq, companyID int64, changePasswordToken string) (int64, error) {

	var res int64

	changePasswordTokenExpirationDate := time.Now().Add(time.Hour * 24 * 7)

	activeFlag := 0
	if userObj.ActiveFlag {
		activeFlag = 1
	}

	query := `
		INSERT
		INTO
			users (email, first_name, last_name, middle_name, position_id, manager_id, phone, active_flag, company_id,
					password_hash, local_iso, is_debug, change_password_token,
					change_password_token_expiration_date, ext_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`

	err := tx.Get(&res, query, userObj.Email, userObj.FirstName, userObj.LastName, userObj.MiddleName, userObj.PositionID,
		userObj.ManagerID, userObj.Phone, activeFlag, companyID, uuid.New(), "rus", false, changePasswordToken,
		changePasswordTokenExpirationDate, userObj.ExtID)

	return res, err
}

func GetRoles(tx *sqlx.Tx, roleIDs []int64, localization_iso string) ([]models.Role, error) {
	var res []models.Role
	query := `
		SELECT
			id,
			ld.value as name
		FROM
			roles r
			LEFT JOIN local_dictionary ld ON ld.key_id = r.name_local_key_id AND ld.language_iso = $2
		WHERE
			id = ANY($1)
	`

	err := tx.Select(&res, query, roleIDs, localization_iso)

	return res, err

}

func AddUserRoles(tx *sqlx.Tx, userID int64, roleIDs []int64) error {
	query := `
		INSERT INTO
			roles_users (user_id, role_id)
		VALUES ($1, $2)
	`

	preparedStmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		_, err := preparedStmt.Exec(userID, roleID)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetLocalizedData(tx *sqlx.Tx, localISO string, ids []int64) ([]models.LocalDictionary, error) {
	var res []models.LocalDictionary

	query := `
		SELECT key_id, language_iso, value
		FROM local_dictionary
		WHERE key_id = ANY($1) AND lower(language_iso) = lower($2)
	`

	err := tx.Select(&res, query, ids, localISO)

	return res, err
}
