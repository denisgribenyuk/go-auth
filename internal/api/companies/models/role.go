package models

type Role struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
	IsSystem    bool   `json:"is_system,omitempty" db:"is_system"`
	Sort        int64  `json:"sort,omitempty" db:"sort"`
	RefName     string `json:"ref_name,omitempty" db:"ref_name"`
}
