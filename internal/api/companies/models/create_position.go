package models

type CreatePositionUrlRequest struct {
	CompanyId int64 `uri:"id" binding:"required"`
}

type CreatePositionBodyRequest struct {
	PositionName string `json:"title" binding:"required,min=1,max=255"`
}
