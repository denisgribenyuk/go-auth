package dbqueries

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetCompanyRoles(companyId int64) (models.CompanyData, error) {
	var company models.CompanyData

	err := postgres.DB.Get(&company, `SELECT c.id,
									c.active_flag,
									c.name,
									COUNT(u.id) AS user_count
								FROM companies c
									LEFT JOIN users u ON c.id = u.company_id
								WHERE c.id = $1
								GROUP BY c.id `, companyId)

	return company, err

}
