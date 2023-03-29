package db_queries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func IsEmailExist(email string) (int, error) {
	var emailCount int
	err := postgres.DB.Get(&emailCount, `
		SELECT count(id) FROM users WHERE email = $1`, email)
	return emailCount, err
}

func IsPhoneExist(phone string) (int, error) {
	var phoneCount int
	err := postgres.DB.Get(&phoneCount, `
		SELECT count(id) FROM users WHERE phone = $1`, phone)
	return phoneCount, err
}

func InsertCompany(tx *sqlx.Tx, comanyName string, isActive bool) (int, error) {
	var companyId int
	err := tx.QueryRowx(`INSERT INTO companies (name, active_flag) VALUES ($1, $2) RETURNING id`, comanyName, isActive).Scan(&companyId)
	return companyId, err
}

func InsertUser(tx *sqlx.Tx, user models.SignUpRequest, companyId int, activeFlag int, isDebug bool, signUpToken string, passwordHash string) (int, error) {
	var userID int
	err := tx.QueryRowx(`INSERT INTO public.users (
		email, 
		first_name, 
		last_name, 
		middle_name, 
		password_hash, 
		local_iso,
		active_flag, 
		company_id, 
		phone, 
		is_debug, 
		sign_up_token) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, user.Email, user.FirstName, user.LastName, user.MiddleName, passwordHash, user.LocalISO, activeFlag, companyId, user.Phone, isDebug, signUpToken).Scan(&userID)
	return userID, err
}

func InsertUserRoles(tx *sqlx.Tx, userId int, roles []int) error {
	for _, role := range roles {
		_, err := tx.Exec(`INSERT INTO roles_users (user_id, role_id) VALUES ($1, $2)`, userId, role)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTranslatedValue(key_id int, language_iso string) (string, error) {
	var value string
	err := postgres.DB.Get(&value, `
		SELECT value FROM local_dictionary WHERE key_id = $1 and language_iso = $2`, key_id, language_iso)
	return value, err
}
