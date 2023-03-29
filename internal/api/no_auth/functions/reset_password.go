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
)

func ResetPassword(c *gin.Context) {
	var request models.ResetPassword
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	//Проверка на существование пользователя с таким email
	_, err := dbqueries.IsEmailExist(request.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error_code": 10001, "error_message": "User with this email not found."})
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	//Генерируем токен
	resetPassToken := CreateRandomString(viper.GetInt("session.tokenLength"))
	//Получаем срок жизни токена
	resetPassTokenExpDate := time.Now().AddDate(0, 0, 1)

	//Получаем пользователя по email
	user, err := dbqueries.GetUser(request.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("getting user error: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	} else if err == sql.ErrNoRows {
		c.String(http.StatusUnauthorized, apierrors.ErrUnauthorized)
		return
	}

	//Записываем токен сброса пользователю
	err = dbqueries.SetChangePasswordToken(resetPassToken, resetPassTokenExpDate, user.ID)
	if err != nil {
		log.Errorf("Change password token error: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}
	//Генерируем уведомление на почту
	message, title, err := createResetPasswordNotification(
		user,
		user.LocalIso,
		resetPassToken,
		request.Email)
	if err != nil {
		log.Errorf("Error in create notification message: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	//Отправляем уведомление на почту
	err = notificator.SendNotification(request.Email, title, message)
	if err != nil {
		log.Errorf("Error in send notification message: %s", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	c.String(http.StatusOK, "OK")

}

func createResetPasswordNotification(user models.DBUser, languageIso string, changePassToken string, email string) (string, string, error) {
	return createNotification(
		user.FirstName,
		user.MiddleName,
		languageIso,
		viper.GetInt("notification.resetPasswordNotificationTileId"),
		viper.GetInt("notification.resetPasswordMessageTileId"),
		"%SUBMIT_URL%",
		"auth/change_password?change_password_token="+changePassToken+"&email"+email,
	)
}
