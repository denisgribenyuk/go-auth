package model

import (
	"fmt"

	"gitlab.assistagro.com/back/back.auth.go/pkg/types"
)

const (
	_                = iota
	RoleAdmin        // 1 - Администратор аккаунта
	RoleUser         // 2 - Пользователь аккаунта
	RoleSuperAdmin   // 3 - Суперадмин
	RoleOperator     // 4 - Оператор справочников
	RoleModerator    // 5 - Модератор справочников
	RoleSalesManager // 6 - Менеджер по продажам
	RoleSupport      // 7 - Специалист техподдержки
	RoleAgronomist   // 8 - Агроном
	RoleDetector     // 9 - Распознавание
)

// User Main user struct
type User struct {
	ID         int64            `json:"id" db:"id"`
	Email      string           `json:"email" db:"email"`
	FirstName  string           `json:"first_name" db:"first_name"`
	LastName   string           `json:"last_name" db:"last_name"`
	MiddleName *string          `json:"middle_name" db:"middle_name"`
	LocalIso   string           `json:"local_iso" db:"local_iso"`
	CompanyID  *int64           `json:"company_id" db:"company_id"`
	Phone      string           `json:"phone" db:"phone"`
	ImageID    *int64           `json:"image_id" db:"image_id"`
	ExtID      *string          `json:"ext_id" db:"ext_id"`
	ManagerID  *int64           `json:"manager_id" db:"manager_id"`
	PositionID *int64           `json:"position_id" db:"position_id"`
	ActiveFlag int64            `json:"active_flag" db:"active_flag"`
	RoleIDS    types.Int64Array `json:"role_ids" db:"role_ids"`
}

func (u User) HasAnyRole(roleIDs ...int64) bool {
	if u.HasRole(RoleSuperAdmin) {
		return true
	}
	for _, roleID := range roleIDs {
		if roleID == 1 {
			return true
		}
		for _, userRoleID := range u.RoleIDS {
			if roleID == userRoleID {
				return true
			}
		}
	}
	return false
}

func (u User) HasAllRoles(roleIDs ...int64) bool {
	fmt.Println("Roles: ", u.RoleIDS)
	if u.HasRole(RoleSuperAdmin) {
		return true
	}
	for _, roleID := range roleIDs {
		found := false
		for _, userRoleID := range u.RoleIDS {
			if roleID == userRoleID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (u User) HasRole(roleID int64) bool {
	for _, userRoleID := range u.RoleIDS {
		if roleID == userRoleID {
			return true
		}
	}
	return false
}
