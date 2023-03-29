package functions

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/service_methods/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/service_methods/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func CheckToken(c *gin.Context) {
	var requestData models.CheckTokenRequest
	err := c.ShouldBindBodyWith(&requestData, binding.JSON)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}

	user, err := dbqueries.GetUserSession(requestData.AccessToken)
	if err != nil && err != sql.ErrNoRows {
		internalError(c, fmt.Sprintf("Getting user session error: %s", err))
		return
	} else if err == sql.ErrNoRows {
		c.String(http.StatusUnauthorized, apierrors.ErrUnauthorized)
		return
	}

	isUserValid, err := user.IsValid()
	if err != nil {
		internalError(c, fmt.Sprintf("User validation error: %s", err))
		return
	}

	if !isUserValid {
		c.String(http.StatusForbidden, apierrors.ErrForbidden)
		return
	}

	res := gin.H{
		"user_id":    user.ID,
		"company_id": user.CompanyID,
		"local_iso":  user.LocalIso,
	}

	if user.HasRole(model.RoleSuperAdmin) {
		c.JSON(http.StatusOK, res)
		return
	}

	if len(requestData.Roles) > 0 {
		if *requestData.RolesRequired && !user.HasAllRoles(requestData.Roles...) ||
			!*requestData.RolesRequired && !user.HasAnyRole(requestData.Roles...) {
			c.String(http.StatusForbidden, apierrors.ErrForbidden)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

func internalError(c *gin.Context, message string) {
	log.Error(message)
	c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
}
