package functions

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"

	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
)

func GetCompanyPositions(context *gin.Context) {
	var requestData models.CompanyURIReq
	err := context.ShouldBindUri(&requestData)
	if err != nil {
		context.String(http.StatusBadRequest, "%", err)
		return
	}
	// Parsing company id
	companyID := requestData.CompanyID
	currentUser := context.MustGet("user").(*model.User)
	isUserIsSuperAdmin := currentUser.HasRole(model.RoleSuperAdmin)
	if err != nil || companyID == 0 {
		context.String(http.StatusBadRequest, apierrors.ErrInvalidCompanyID)
		return
	}

	//Проверка на существование компании
	_, err = dbqueries.GetCompanyData(companyID)
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
	isUserInCompany := currentUser.CompanyID != nil && *currentUser.CompanyID == companyID
	if !isUserInCompany && !isUserIsSuperAdmin {
		context.String(http.StatusForbidden, apierrors.ErrUserNotInCompany)
		return
	}

	//Получение должностей
	positions, err := dbqueries.GetCompanyPositions(companyID)
	if err != nil {
		// Other db errors
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	var res = map[string][]models.Position{"positions": positions}
	context.JSON(http.StatusOK, res)
}
