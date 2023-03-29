package dbqueries

import (
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func CreatePositionInCompany(companyID int64, positionName string) (int64, error) {
	var id int64
	err := postgres.DB.QueryRow(
		`INSERT INTO company_positions (id, company_id, name)
				VALUES (
						(SELECT
							CASE
								WHEN
									max(id) IS NULL
								THEN
									1
								ELSE
									max(id)+1
								END
						FROM
							company_positions),
				        $1, $2) RETURNING id`, companyID, positionName).Scan(&id)
	return id, err
}

func GetPositionByName(positionName string, companyId int64) (string, error) {
	var title string
	err := postgres.DB.Get(&title,
		`SELECT name FROM company_positions WHERE name=$1 AND company_id=$2`, positionName, companyId)
	return title, err
}
