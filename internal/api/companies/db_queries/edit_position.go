package dbqueries

import (
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func EditPositionInCompany(PositionTitle string, positionId int64) error {
	_, err := postgres.DB.Exec(
		`UPDATE company_positions
				SET name = $1
				 WHERE id = $2`, PositionTitle, positionId)
	return err
}

func GetPositionByIdAndCompany(positionId int64, companyId int64) (string, error) {
	var res string
	err := postgres.DB.Get(&res,
		`SELECT name FROM company_positions WHERE id=$1 AND company_id=$2`, positionId, companyId)
	return res, err
}
