package model

type CompanyPosition struct {
	ID        int64  `json:"id" db:"id"`
	CompanyID int64  `json:"company_id" db:"company_id"`
	Name      string `json:"name" db:"name"`
}
