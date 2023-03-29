package functions

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/notificator"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func ResendEmail(c *gin.Context) {
	time.Sleep(1 * time.Second)

	var requestData models.ResendEmail
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.String(http.StatusBadRequest, "%", err)
		return
	}

	user, err := dbqueries.GetUserF(requestData.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("getting user error: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	} else if err == sql.ErrNoRows {
		c.String(http.StatusUnauthorized, "User with this email not found.")
		return
	}

	if !user.SignUPToken.Valid {
		c.String(http.StatusBadRequest, "This email already verified.")
		return
	}

	signUpToken := CreateRandomString(viper.GetInt("session.tokenLength"))
	txDB, err := postgres.DB.Beginx()
	if err != nil {
		log.Errorf("Error in open transaction: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}
	defer txDB.Rollback()

	err = dbqueries.UpdateUsersignUpToken(txDB, user.ID, signUpToken)
	if err != nil {
		log.Errorf("Error in default roles_users isert: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	message, title, err := createSignUpNotification(user.FirstName,
		user.MiddleName,
		user.LocalIso,
		signUpToken)
	if err != nil {
		log.Errorf("Error in create notification message: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	err = notificator.SendNotification(user.Email, title, message)
	if err != nil {
		log.Errorf("Error in send notification message: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	txDB.Commit()

	c.String(http.StatusOK, "OK")
}
