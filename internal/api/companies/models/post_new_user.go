package models

type PostNewUserReq struct {
	ExtID        *string  `json:"ext_id" binding:"omitempty,min=1,max=255"`
	FirstName    string   `json:"first_name" binding:"required,min=1,max=255"`
	LastName     string   `json:"last_name" binding:"required,min=1,max=255"`
	MiddleName   *string  `json:"middle_name" binding:"omitempty,min=1,max=255"`
	Email        string   `json:"email" binding:"required,email,min=1,max=255"`
	Phone        string   `json:"phone" binding:"required,phone"`
	PositionID   *int64   `json:"position_id" binding:"omitempty,min=1"`
	ManagerID    *int64   `json:"manager_id" binding:"omitempty,min=1"`
	StructureIDs []int64  `json:"structure_ids" binding:"omitempty,dive,min=1"`
	RoleIDs      []int64  `json:"role_ids" binding:"omitempty,dive,min=1"`
	FieldGUIDs   []string `json:"field_guids" binding:"omitempty,dive,uuid"`
	ActiveFlag   bool     `json:"active_flag" binding:"required"`
}
