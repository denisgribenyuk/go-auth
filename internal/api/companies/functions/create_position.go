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

func CreatePosition(context *gin.Context) {
	var requestDataUrl models.CreatePositionUrlRequest
	var requestDataBody models.CreatePositionBodyRequest
	err := context.ShouldBindUri(&requestDataUrl)
	if err != nil {
		context.String(http.StatusBadRequest, "%", err)
		return
	}
	err = context.ShouldBindJSON(&requestDataBody)
	if err != nil {
		context.String(http.StatusBadRequest, "%", err)
		return
	}
	// Get user from context
	currentUser := context.MustGet("user").(*model.User)
	// Parsing body
	companyId := requestDataUrl.CompanyId
	positionName := requestDataBody.PositionName

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

	name, err := dbqueries.GetPositionByName(positionName, companyId)
	if err != nil && err != sql.ErrNoRows {
		context.String(http.StatusInternalServerError, err.Error())
		return
	} else if err == nil && name != "" {
		context.String(http.StatusConflict, apierrors.ErrPositionsConflict)
		return
	}

	newPositionId, err := dbqueries.CreatePositionInCompany(companyId, positionName)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	var res = map[string]int64{"id": newPositionId}
	context.JSON(http.StatusOK, res)
}
