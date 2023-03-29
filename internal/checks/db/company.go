package checks

import (
	"database/sql"

	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func IsCompanyExists(companyID int64) (bool, error) {
	var res interface{}
	err := postgres.DB.Get(
		&res,
		`
		SELECT id
					FROM company
					WHERE id = $1
		`,
		companyID)

	return err == nil, checkNoRows(err)

}

func IsUserInCompany(userId int64, companyId int64) (bool, error) {
	var res interface{}
	err := postgres.DB.Get(
		&res,
		`SELECT id
				FROM users
				WHERE id = $1 AND company_id = $2`,
		userId, companyId)

	return err == nil, checkNoRows(err)
}

func checkNoRows(err error) error {
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}
