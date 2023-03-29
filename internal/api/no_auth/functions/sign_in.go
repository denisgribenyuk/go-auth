package functions

import (
	"crypto/hmac"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	internal_models "gitlab.assistagro.com/back/back.auth.go/internal/models"
)

func SignIN(c *gin.Context) {
	var requestData models.SignINRequest
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.String(http.StatusBadRequest, "%", err)
		return
	}

	user, err := db_queries.GetUser(requestData.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("getting user error: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	} else if err == sql.ErrNoRows {
		c.String(http.StatusUnauthorized, apierrors.ErrUnauthorized)
		return
	}

	if ok, err := verifyPassword(requestData.Password, user.PasswordHash); !ok {
		if err != nil {
			log.Errorf("Verify password error: %s", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		} else {
			c.String(http.StatusUnauthorized, "Incorrect email or password.")
			return
		}
	}

	if user.SignUPToken.Valid {
		c.String(http.StatusTooEarly, "Email not verified.")
		return
	}

	isUserValid, err := user.IsValid()
	if err != nil {
		log.Errorf("User validation error: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	if isUserValid {
		session, err := createSession(user)
		if err != nil {
			log.Errorf("Creating session error: %s", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":             session.AccessToken,
			"access_token_expires_in":  viper.GetFloat64("session.accessTTL"),
			"refresh_token":            session.RefreshToken,
			"refresh_token_expires_in": viper.GetFloat64("session.refreshTTL"),
		})
		return
	}

	c.String(http.StatusForbidden, apierrors.ErrForbidden)
}

func verifyPassword(password string, encoded string) (bool, error) {
	s := strings.Split(encoded, "$")
	if len(s) != 3 {
		return false, errors.New("hashed password components mismatch")
	}
	infoComponents := strings.Split(s[0], ":")
	algorithm, algorithm2, iterations := infoComponents[0], infoComponents[1], infoComponents[2]
	if algorithm != "pbkdf2" || algorithm2 != "sha256" {
		return false, errors.New("algorithm mismatch")
	}
	i, err := strconv.Atoi(iterations)
	if err != nil {
		return false, errors.New("unreadable component in hashed password")
	}
	salt := s[1]

	newEncoded, err := PasswordEncode(password, salt, i)
	if err != nil {
		return false, err
	}

	return hmac.Equal([]byte(newEncoded), []byte(encoded)), nil

}

func createSession(user models.DBUser) (internal_models.Session, error) {
	accessTTL := time.Duration(viper.GetInt("session.accessTTL")) * time.Second
	refreshTTL := time.Duration(viper.GetInt("session.refreshTTL")) * time.Second
	tokenLen := viper.GetInt("session.tokenLength")

	session := internal_models.Session{
		UserID:                sql.NullInt64{Int64: user.ID, Valid: true},
		AccessToken:           CreateRandomString(tokenLen),
		RefreshToken:          CreateRandomString(tokenLen),
		AccessExpirationDate:  time.Now().UTC().Add(accessTTL),
		RefreshExpirationDate: time.Now().UTC().Add(refreshTTL),
	}

	err := db_queries.CreateSession(session)
	return session, err
}
