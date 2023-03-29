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

func EditPosition(context *gin.Context) {
	var requestDataUrl models.EditPositionUrlRequest
	var requestDataBody models.EditPositionBodyRequest
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
	positionId := requestDataUrl.PositionId
	positionName := requestDataBody.PositionName

	isUserSuperAdmin := currentUser.HasRole(model.RoleSuperAdmin)

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
	if !isUserInCompany && !isUserSuperAdmin {
		context.String(http.StatusForbidden, apierrors.ErrUserNotInCompany)
		return
	}

	//Проверка, что новое имя должности не занято в компании
	name, err := dbqueries.GetPositionByName(positionName, companyId)
	if err != nil && err != sql.ErrNoRows {
		context.String(http.StatusInternalServerError, err.Error())
		return
	} else if err == nil && name != positionName {
		context.String(http.StatusConflict, apierrors.ErrPositionsConflict)
		return
	}

	//Проверка, что должность с таким id есть в компании
	_, err = dbqueries.GetPositionByIdAndCompany(positionId, companyId)
	if err != nil && err != sql.ErrNoRows {
		context.String(http.StatusInternalServerError, err.Error())
		return
	} else if err == sql.ErrNoRows {
		context.String(http.StatusNotFound, apierrors.ErrPositionNotInCompany)
		return
	}

	err = dbqueries.EditPositionInCompany(positionName, positionId)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	context.JSON(http.StatusOK, "")
}
