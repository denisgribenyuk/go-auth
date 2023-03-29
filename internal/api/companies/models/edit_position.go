package models

type EditPositionUrlRequest struct {
	CompanyId  int64 `uri:"id" binding:"required,min:1"`
	PositionId int64 `uri:"position_id" binding:"required,min:1"`
}

type EditPositionBodyRequest struct {
	PositionName string `json:"title" binding:"required,min=1,max=255"`
}
