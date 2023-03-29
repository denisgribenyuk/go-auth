package models

import "gitlab.assistagro.com/back/back.auth.go/pkg/model"

type CompanyUserURIReq struct {
	CompanyID int64 `uri:"id" binding:"required,min=1"`
	UserID    int64 `uri:"user_id" binding:"required,min=1"`
}
type GetUserImageQueryReq struct {
	IsThumbnail *int `form:"is_thumbnail" binding:"required,oneof=0 1"`
}

type UserPhoto struct {
	Image string `json:"image" db:"image" binding:"required,max=10000000"`
	Name  string `json:"name" db:"name" binding:"required,max=512"`
}

type DBUserImage struct {
	model.User
	ImagePath string `json:"image_path" db:"image_path"`
}
