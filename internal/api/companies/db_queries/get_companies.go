package dbqueries

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetCompanyByAPIKey(apiKey string) (models.CompanyData, error) {
	var company models.CompanyData

	err := postgres.DB.Get(&company, `
		SELECT  c.id,
				c.active_flag,
				c.name,
				COUNT(u.id) AS user_count
		FROM companies c
				LEFT JOIN users u ON c.id = u.company_id
		WHERE c.api_key = $1
		GROUP BY c.id `, apiKey)
	return company, err
}
