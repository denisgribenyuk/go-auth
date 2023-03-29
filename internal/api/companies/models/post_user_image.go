package models

type PostUserImageQueryReq struct {
	Base64 *string `form:"base64" binding:"required,max=15000000"`
}
