package dbqueries

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetCompanyPositions(companyId int64) ([]models.Position, error) {
	positions := []models.Position{}

	err := postgres.DB.Select(
		&positions,
		`SELECT id, "name"
		FROM company_positions
		WHERE company_id = $1
		ORDER BY name`,
		companyId)
	return positions, err
}
