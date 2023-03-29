package functions

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func GetCompanyRoles(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	// Parsing company id
	companyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || companyID == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company id"})
		c.String(http.StatusBadRequest, apierrors.ErrInvalidCompanyID)
		return
	}

	// Проверяем существование компании
	company, err := dbqueries.GetCompanyData(companyID)
	if err != nil {
		// If empty result
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, apierrors.ErrCompanyNotFound)
			return
		}
		// Other db errors
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Проверяем, принадлежит ли пользователь компании
	if !user.HasRole(model.RoleSuperAdmin) {
		if company.ID != *user.CompanyID {
			c.String(http.StatusForbidden, apierrors.ErrUserNotInCompany)
			return
		}
	}

	// Составляем список ролей
	roles := []models.Role{
		{
			ID:   1,
			Name: "Администратор",
		},
		{
			ID:   2,
			Name: "Агроном",
		},
		{
			ID:   3,
			Name: "Суперадминистратор",
		},
		{
			ID:   8,
			Name: "Скаут",
		},
		{
			ID:   9,
			Name: "Исследователь",
		},
	}

	c.JSON(http.StatusOK, roles)
}
