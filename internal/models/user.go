package models

import (
	"database/sql"

	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

type User struct {
	model.User
}

func (user *User) IsValid() (bool, error) {
	if user.CompanyID == nil {
		return user.ActiveFlag == 1, nil
	}

	var userCompany company
	query := "SELECT id, active_flag FROM companies WHERE id = $1"
	err := postgres.DB.Get(&userCompany, query, user.CompanyID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows {
		return user.ActiveFlag == 1, nil
	} else {
		return user.ActiveFlag == 1 && userCompany.ActiveFlag, nil
	}
}

type company struct {
	ID         int  `db:"id"`
	ActiveFlag bool `db:"active_flag"`
}
