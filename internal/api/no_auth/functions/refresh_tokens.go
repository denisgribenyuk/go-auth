package functions

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	internalmodels "gitlab.assistagro.com/back/back.auth.go/internal/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func RefreshTokens(c *gin.Context) {
	var requestData models.RefreshTokensRequest
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}

	tx, err := postgres.DB.Beginx()
	if err != nil {
		NewInternalServerError(c, fmt.Sprintf("Transaction opening error: %s", err))
		return
	}
	defer tx.Rollback()

	session, err := db_queries.GetSession(tx, requestData.RefreshToken)
	if err != nil {
		if err != sql.ErrNoRows {
			NewInternalServerError(c, fmt.Sprintf("Getting session error: %s", err))
			return
		} else {
			c.String(http.StatusUnauthorized, apierrors.ErrUnauthorized)
			return
		}
	}

	user, err := db_queries.GetUserByID(session.UserID.Int64)
	if err != nil {
		NewInternalServerError(c, fmt.Sprintf("Getting user error: %s", err))
		return
	}

	isUserValid, err := user.IsValid()
	if err != nil {
		NewInternalServerError(c, fmt.Sprintf("User validation error: %s", err))
		return
	}

	if !isUserValid {
		c.String(http.StatusForbidden, apierrors.ErrForbidden)
		return
	}

	err = updateSession(tx, session)
	if err != nil {
		NewInternalServerError(c, fmt.Sprintf("Updating session error: %s", err))
		return
	}

	err = tx.Commit()
	if err != nil {
		message := fmt.Sprintf("Transaction commit error: %s", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			message = fmt.Sprintf("Transaction commit error. Message: %s. Details: %s.", pgErr.Message, pgErr.Detail)
		}
		NewInternalServerError(c, message)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":             session.AccessToken,
		"access_token_expires_in":  viper.GetFloat64("session.accessTTL"),
		"refresh_token":            session.RefreshToken,
		"refresh_token_expires_in": viper.GetFloat64("session.refreshTTL"),
	})
}

func updateSession(tx *sqlx.Tx, session *internalmodels.Session) error {
	accessTTL := time.Duration(viper.GetInt("session.accessTTL")) * time.Second
	refreshTTL := time.Duration(viper.GetInt("session.refreshTTL")) * time.Second
	tokenLen := viper.GetInt("session.tokenLength")

	session.AccessToken = CreateRandomString(tokenLen)
	session.AccessExpirationDate = time.Now().UTC().Add(accessTTL)
	session.RefreshToken = CreateRandomString(tokenLen)
	session.RefreshExpirationDate = time.Now().UTC().Add(refreshTTL)

	err := db_queries.UpdateSession(tx, session)
	return err
}
