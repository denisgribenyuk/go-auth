package functions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func GetCompanyData(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(*model.User)

	// Parsing company id
	companyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || companyID == 0 {
		c.String(http.StatusBadRequest, apierrors.ErrInvalidCompanyID)
		return
	}

	// Get company data
	company, err := dbqueries.GetCompanyData(companyID)
	if err != nil {
		// If empty result
		if err == pgx.ErrNoRows {
			c.String(http.StatusNotFound, apierrors.ErrCompanyNotFound)
			return
		}
		// Other db-related errors
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// if user is not superadmin, check if user is in company
	if !user.HasRole(model.RoleSuperAdmin) {
		if company.ID != *user.CompanyID {
			c.String(http.StatusForbidden, apierrors.ErrUserNotInCompany)
			return
		}
	}

	c.JSON(http.StatusOK, company)
}
