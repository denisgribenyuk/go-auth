package functions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func GetCurrentUser(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"middle_name": user.MiddleName,
		"phone":       user.Phone,
		"company_id":  user.CompanyID,
		"role_ids":    user.RoleIDS,
	})
}
