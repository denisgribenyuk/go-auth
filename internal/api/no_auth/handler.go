package no_auth

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/functions"
	"gitlab.assistagro.com/back/back.auth.go/internal/transport/rest"
)

func Register(r *rest.Router) {
	serviceGroup := r.Group("/")

	serviceGroup.POST("/sign_up", functions.SignUp)
	serviceGroup.POST("/sign_in", functions.SignIN)
	serviceGroup.POST("/refresh_tokens", functions.RefreshTokens)
	serviceGroup.POST("/resend_email", functions.ResendEmail)
	serviceGroup.POST("/reset_password", functions.ResetPassword)
}
