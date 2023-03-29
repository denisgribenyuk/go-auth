package functions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	log "github.com/sirupsen/logrus"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/tools"
)

func GetCompanyUsers(c *gin.Context) {
	allowedParams := []string{"limit", "offset", "sort_column", "sort_direction", "user_id", "ext_id", "user_name",
		"position_id", "email", "phone", "role_id", "manager_id", "structure_id", "field_guid", "active_flag"}
	for k := range c.Request.URL.Query() {
		if !tools.Contains(allowedParams, k) {
			c.String(http.StatusBadRequest, apierrors.ErrExtraArgsNotAllowed)
			return
		}
	}

	var uriData models.GetCompanyUsersRequestURI
	err := c.ShouldBindUri(&uriData)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}

	var queryData models.GetCompanyUsersRequestQuery
	err = c.ShouldBindQuery(&queryData)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}

	users, totalCount, err := dbqueries.GetCompanyUsers(uriData.CompanyID, queryData)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Errorf("Getting users error. Code:%s. Message: %s. Details: %s.",
				pgErr.Code, pgErr.Message, pgErr.Detail)
		} else {
			log.Errorf("Getting users error: %s", err)
		}
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	limit := 0
	if queryData.Limit != nil {
		limit = *queryData.Limit
	}
	res := models.GetCompanyUsersResponse{
		Users: users,
		Paging: models.DBPaging{
			Limit:      limit,
			Offset:     queryData.Offset,
			TotalCount: totalCount,
		},
	}

	c.JSON(http.StatusOK, res)
}
