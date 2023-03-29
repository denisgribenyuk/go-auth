package models

import (
	"database/sql"
	"time"
)

type Session struct {
	ID                    sql.NullInt64 `db:"id"`
	UserID                sql.NullInt64 `db:"user_id"`
	AccessToken           string        `db:"access_token"`
	RefreshToken          string        `db:"refresh_token"`
	AccessExpirationDate  time.Time     `db:"access_expiration_date"`
	RefreshExpirationDate time.Time     `db:"refresh_expiration_date"`
}
