package authmiddleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/user/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithError(code, fmt.Errorf("%v", message))
	c.String(code, "%v", message)
}

func AuthorizedUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if api_key is in context
		if _, ok := c.Get("api_key"); ok {
			c.Next()
			return
		}

		token := c.Request.Header.Get("X-Token")
		if token == "" {
			respondWithError(c, 400, "Token must be in request headers.")
			return
		}

		// Get user from db
		user, err := dbqueries.GetUserByToken(token)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, 401, "The session does not exist or has expired.")
				return
			}
			respondWithError(c, 500, err)
			return
		}

		// Put user to context
		c.Set("user", user)
		c.Next()
	}
}

func RolesAcceptedMiddleware(roles ...int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*model.User)
		if !user.HasAnyRole(roles...) {
			respondWithError(c, 403, "Access denied")
			return
		}
		c.Next()
	}
}

func RolesRequiredMiddleware(roles ...int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*model.User)
		if !user.HasAllRoles(roles...) {
			respondWithError(c, 403, "Access denied")
			return
		}
		c.Next()
	}
}

func AuthorizedServiceRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.Request.Header.Get("X-Service-Token"); token != os.Getenv("SERVICE_SECRET_KEY") {
			respondWithError(c, http.StatusNotFound, apierrors.ErrURLNotFound)
			return
		}
		c.Next()
	}
}

func UserOrServiceRequired(roles ...int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.Request.Header.Get("X-Token")
		serviceToken := c.Request.Header.Get("X-Service-Token")
		if userToken != "" {
			AuthorizedUserMiddleware()(c)
			RolesAcceptedMiddleware(roles...)(c)
		} else if serviceToken != "" {
			AuthorizedServiceRequired()(c)
		} else {
			respondWithError(c, http.StatusNotFound, apierrors.ErrURLNotFound)
		}
	}
}
