package dbqueries

import (
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func ClearPositionInUsers(positionID int64, companyID int64) error {
	_, err := postgres.DB.Exec(
		`UPDATE users
		SET position_id = null
		WHERE company_id = $1
			AND position_id = $2`, companyID, positionID)
	return err
}

func DeletePositionInCompany(positionID int64, companyID int64) error {
	_, err := postgres.DB.Exec(
		`DELETE FROM company_positions
		WHERE company_id = $1
			AND id = $2`, companyID, positionID)
	return err
}
