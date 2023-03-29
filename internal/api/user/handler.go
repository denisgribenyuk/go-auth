package user

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/user/functions"
	"gitlab.assistagro.com/back/back.auth.go/internal/middleware/authmiddleware"
	"gitlab.assistagro.com/back/back.auth.go/internal/transport/rest"
)

func Register(r *rest.Router) {
	userGroup := r.Group("/", authmiddleware.AuthorizedUserMiddleware())

	userGroup.GET("/current_user", functions.GetCurrentUser)
}
