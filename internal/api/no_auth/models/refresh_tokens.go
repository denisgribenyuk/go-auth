package models

type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,min=80,max=80"`
}
