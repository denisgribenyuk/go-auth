package service_methods

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/service_methods/functions"
	"gitlab.assistagro.com/back/back.auth.go/internal/middleware/authmiddleware"
	"gitlab.assistagro.com/back/back.auth.go/internal/transport/rest"
)

func Register(r *rest.Router) {
	serviceGroup := r.Group("/service", authmiddleware.AuthorizedServiceRequired())

	serviceGroup.POST("/check_token", functions.CheckToken)
}
