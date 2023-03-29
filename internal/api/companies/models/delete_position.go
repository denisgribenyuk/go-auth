package models

type DeletePositionRequest struct {
	CompanyId  int64 `uri:"id" binding:"required"`
	PositionId int64 `uri:"position_id" binding:"required"`
}
