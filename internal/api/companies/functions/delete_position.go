package functions

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func DeletePosition(context *gin.Context) {
	var requestData models.DeletePositionRequest
	err := context.ShouldBindUri(&requestData)
	if err != nil {
		context.String(http.StatusBadRequest, "%", err)
		return
	}
	// Get user from context
	currentUser := context.MustGet("user").(*model.User)
	// Parsing body
	companyId := requestData.CompanyId
	positionId := requestData.PositionId

	isUserIsSuperAdmin := currentUser.HasRole(model.RoleSuperAdmin)

	//Проверка на существование компании
	_, err = dbqueries.GetCompanyData(companyId)
	if err != nil {
		// If empty result
		if err == sql.ErrNoRows {
			context.String(http.StatusNotFound, apierrors.ErrCompanyNotFound)
			return
		} else {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	//Проверка принадлежности пользователя компании
	isUserInCompany := currentUser.CompanyID != nil && *currentUser.CompanyID == companyId
	if !isUserInCompany && !isUserIsSuperAdmin {
		context.String(http.StatusForbidden, apierrors.ErrUserNotInCompany)
		return
	}

	err = dbqueries.ClearPositionInUsers(positionId, companyId)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = dbqueries.DeletePositionInCompany(positionId, companyId)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	context.JSON(http.StatusOK, "")
}
