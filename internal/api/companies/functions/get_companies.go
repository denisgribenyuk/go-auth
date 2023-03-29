package functions

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
)

func GetCompanies(c *gin.Context) {
	// Get api_key from context
	apiKey := c.Query("api_key")
	err := binding.Validator.Engine().(*validator.Validate).Var(apiKey, "omitempty,min=40,max=40")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get company
	company, err := dbqueries.GetCompanyByAPIKey(apiKey)
	if err != nil {
		// If empty result
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, apierrors.ErrCompanyNotFound)
			return
		}
		// Other db-related errors
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int64{"company_id": company.ID})
}
