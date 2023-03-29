package models

type UserStructure struct {
	UserID      int64 `json:"user_id" db:"user_id"`
	StructureID int64 `json:"structure_id" db:"structure_id"`
}
