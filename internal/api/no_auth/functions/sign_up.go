package functions

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/notificator"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func SignUp(c *gin.Context) {

	var request models.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	request.LocalISO = strings.ToLower(request.LocalISO)
	// TODO: remove
	fmt.Println(request)

	emailCount, err := dbqueries.IsEmailExist(request.Email)
	if err != nil || emailCount != 0 {
		c.JSON(http.StatusConflict, gin.H{"error_code": 10001, "error_message": "The email is already in use."})
		return
	}

	phoneCount, err := dbqueries.IsPhoneExist(request.Phone)
	if err != nil || phoneCount != 0 {
		c.JSON(http.StatusConflict, gin.H{"error_code": 10002, "error_message": "The phone number is already in use."})
		return
	}

	txDB, err := postgres.DB.Beginx()
	if err != nil {
		log.Errorf("Error in open transaction: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}
	defer txDB.Rollback()

	isActiveCompany := true
	companyId, err := dbqueries.InsertCompany(txDB, request.CompanyName, isActiveCompany)
	if err != nil {
		log.Errorf("Error in company isert: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// TODO: remove
	fmt.Println("Company id = ", companyId)

	isActiveUser := 1
	isDebugUser := false
	signUpToken := CreateRandomString(viper.GetInt("session.tokenLength"))

	passwordHash, err := PasswordEncode(request.Password, "", 0)
	if err != nil {
		log.Errorf("Error in PasswordEncode: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
	}

	userID, err := dbqueries.InsertUser(txDB, request, companyId, isActiveUser, isDebugUser, signUpToken, passwordHash)
	if err != nil {
		log.Errorf("Error in user isert: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}
	// TODO: remove
	fmt.Println(userID)

	defaultRoles := []int{model.RoleAdmin, model.RoleUser}
	err = dbqueries.InsertUserRoles(txDB, userID, defaultRoles)
	if err != nil {
		log.Errorf("Error in default roles_users isert: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	txDB.Commit()

	message, title, err := createSignUpNotification(request.FirstName,
		request.MiddleName,
		request.LocalISO,
		signUpToken)
	if err != nil {
		log.Errorf("Error in create notification message: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	err = notificator.SendNotification(request.Email, title, message)
	if err != nil {
		log.Errorf("Error in send notification message: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	c.String(http.StatusOK, "OK")
}

func createSignUpNotification(firstName string, middleName *string, language_iso string, signUpToken string) (string, string, error) {
	return createNotification(firstName,
		middleName,
		language_iso,
		viper.GetInt("notification.signUpNotificationTitleId"),
		viper.GetInt("notification.signUpNotificationMessageId"),
		"%SUBMIT_URL%",
		"/auth/confirm_sign_up?sign_up_token="+signUpToken)
}

func createNotification(firstName string, middleName *string, language_iso string, titelId int, messageId int, submitUrlPattern string, URLPath string) (string, string, error) {
	title, err := dbqueries.GetTranslatedValue(titelId, language_iso)
	if err != nil {
		return "", "", err
	}

	message, err := dbqueries.GetTranslatedValue(messageId, language_iso)
	if err != nil {
		return "", "", err
	}

	fullName := firstName
	if middleName != nil {
		fullName += " " + *middleName
	}

	message = strings.Replace(message, "%FULL_NAME%", fullName, -1)
	message = strings.Replace(message, submitUrlPattern, os.Getenv("FRONT_URL")+URLPath, -1)
	message = strings.Replace(message, "%FRONT_URL%", os.Getenv("FRONT_URL"), -1)
	message = strings.Replace(message, "%GOOGLE_PLAY_APP_LINK%", os.Getenv("GOOGLE_PLAY_APP_LINK"), -1)
	message = strings.Replace(message, "%CURRENT_YEAR%", strconv.Itoa(time.Now().Year()), -1)
	message = strings.Replace(message, "%TITLE%", title, -1)

	return message, title, nil
}
