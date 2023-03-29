package companies

import (
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/functions"
	"gitlab.assistagro.com/back/back.auth.go/internal/middleware/authmiddleware"
	"gitlab.assistagro.com/back/back.auth.go/internal/transport/rest"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func Register(r *rest.Router) {
	companiesUserAuthGroup := r.Group("/companies", authmiddleware.AuthorizedUserMiddleware())

	companiesUserAuthGroup.POST("/:id/users",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleSuperAdmin),
		functions.PostNewUser)

	companiesUserAuthGroup.GET("/:id/users/:user_id/image",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleUser, model.RoleSuperAdmin),
		functions.GetUserImage)

	companiesUserAuthGroup.POST("/:id/users/:user_id/image",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleUser, model.RoleSuperAdmin),
		functions.PostUserImage)

	companiesUserAuthGroup.PATCH("/:id/users/:user_id",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleUser, model.RoleSuperAdmin),
		functions.PatchUserImage)

	companiesUserAuthGroup.GET("/:id",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleUser, model.RoleSuperAdmin),
		functions.GetCompanyData)

	companiesUserAuthGroup.GET("/:id/roles",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleSuperAdmin),
		functions.GetCompanyRoles)

	companiesUserAuthGroup.GET("/:id/positions",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleSuperAdmin),
		functions.GetCompanyPositions)

	companiesUserAuthGroup.DELETE("/:id/positions/:position_id",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleSuperAdmin),
		functions.DeletePosition)

	companiesUserAuthGroup.POST("/:id/positions/",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleSuperAdmin),
		functions.CreatePosition)

	companiesUserAuthGroup.PATCH("/:id/positions/:position_id",
		authmiddleware.RolesAcceptedMiddleware(model.RoleAdmin, model.RoleSuperAdmin),
		functions.EditPosition)

	companiesServiceAuthGroup := r.Group("/companies",
		authmiddleware.AuthorizedServiceRequired())
	companiesServiceAuthGroup.GET("", functions.GetCompanies)

	companiesServiceUserAuthGroup := r.Group("/companies")
	companiesServiceUserAuthGroup.GET("/:id/users",
		authmiddleware.UserOrServiceRequired(model.RoleAdmin, model.RoleUser, model.RoleSuperAdmin),
		functions.GetCompanyUsers)
}
