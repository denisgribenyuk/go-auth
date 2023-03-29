package models

import (
	"github.com/lib/pq"
)

type GetCompanyUsersRequestURI struct {
	CompanyID int `uri:"id" binding:"required"`
}

type GetCompanyUsersRequestQuery struct {
	Limit         *int     `form:"limit" binding:"omitempty,gt=0,lt=101"`
	Offset        int      `form:"offset" binding:"omitempty,gte=0"`
	SortColumn    *string  `form:"sort_column" binding:"omitempty,oneof=last_name position_name email phone ext_id manager_last_name"`
	SortDirection *string  `form:"sort_direction" binding:"omitempty,oneof=asc ASC desc DESC"`
	UserIDs       []int    `form:"user_id" binding:"omitempty,lte=10,dive"`
	ExtIDs        []string `form:"ext_id" binding:"omitempty,lte=10,dive,max=255"`
	UserName      *string  `form:"user_name" binding:"omitempty,min=1,max=255"`
	PositionIDs   []string `form:"position_id" binding:"omitempty"`
	Email         *string  `form:"email" binding:"omitempty,min=1,max=255"`
	Phone         *string  `form:"phone" binding:"omitempty,min=1,max=255"`
	RoleIDs       []string `form:"role_id" binding:"omitempty"`
	ManagerIDs    []string `form:"manager_id" binding:"omitempty"`
	StructureIDs  []string `form:"structure_id" binding:"omitempty"`
	FieldGUIDs    []string `form:"field_guid" binding:"omitempty"`
	ActiveFlag    *bool    `form:"active_flag" binding:"omitempty"`
}

type DBCompanyUsers struct {
	ID              int            `db:"id" json:"id"`
	ExtID           *string        `db:"ext_id" json:"ext_id"`
	UserName        string         `db:"user_name" json:"user_name"`
	FirstName       string         `db:"first_name" json:"first_name"`
	LastName        string         `db:"last_name" json:"last_name"`
	MiddleName      *string        `db:"middle_name" json:"middle_name"`
	Email           string         `db:"email" json:"email"`
	Phone           string         `db:"phone" json:"phone"`
	CompanyID       int            `db:"company_id" json:"company_id"`
	ImageName       *string        `db:"image_name" json:"image_name"`
	PositionID      *int           `db:"position_id" json:"position_id"`
	PositionName    *string        `db:"position_name" json:"position_name"`
	ManagerID       *int           `db:"manager_id" json:"manager_id"`
	ManagerLastName *string        `db:"manager_last_name" json:"manager_last_name"`
	StructureIDs    pq.Int64Array  `db:"structure_ids" json:"structure_ids"`
	RoleIDs         pq.Int64Array  `db:"role_ids" json:"role_ids"`
	FieldGUIDs      pq.StringArray `db:"field_guids" json:"field_guids"`
	ActiveFlag      bool           `db:"active_flag" json:"active_flag"`
}

type DBPaging struct {
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	TotalCount int `json:"total_count"`
}

type GetCompanyUsersResponse struct {
	Users  []DBCompanyUsers `json:"users"`
	Paging DBPaging         `json:"paging"`
}
