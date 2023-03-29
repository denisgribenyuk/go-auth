package models

import "github.com/google/uuid"

type UserField struct {
	UserID    int64     `json:"user_id" db:"user_id"`
	FieldGUID uuid.UUID `json:"field_guid" db:"field_guid"`
}
