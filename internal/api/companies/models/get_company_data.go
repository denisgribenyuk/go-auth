package models

type CompanyData struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	UserCount  string `json:"user_count" db:"user_count"`
	ActiveFlag string `json:"active_flag" db:"active_flag"`
}
