package models

type CheckTokenRequest struct {
	AccessToken   string  `json:"access_token" binding:"required,min=80,max=80"`
	Roles         []int64 `json:"roles" binding:"required"`
	RolesRequired *bool   `json:"roles_required" binding:"required"`
}
