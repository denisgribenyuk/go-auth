package models

type Position struct {
	ID    int64  `json:"id" db:"id"`
	Title string `json:"title" db:"name"`
}

type UserInCompany struct {
	ID int64 `json:"id" db:"id"`
}

type CompanyURIReq struct {
	CompanyID int64 `uri:"id" binding:"required"`
}
