package functions

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	companyqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	userqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/db_queries"
	noauthfunctions "gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/functions"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/backweb/bwclient"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/notificator"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func PostNewUser(c *gin.Context) {

	var newUserURIReq models.CompanyURIReq
	var newUserReq models.PostNewUserReq

	// Обрабатываем параметры пути
	err := c.ShouldBindUri(&newUserURIReq)
	if err != nil {
		log.Errorf("Error binding uri: %w", err)
		c.String(http.StatusNotFound, apierrors.ErrURLNotFound)
		return
	}

	// Обрабатываем параметры запроса
	err = c.ShouldBindJSON(&newUserReq)
	if err != nil {
		log.Errorf("Error binding json: %w", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Вытаскиваем пользователя
	user := c.MustGet("user").(*model.User)
	user.LocalIso = strings.ToLower(user.LocalIso)

	// Для роли 1 проверяем доступ к компании, для которой создается пользователь. Для роли 3 - можно создавать пользователей для всех компаний
	if !user.HasRole(model.RoleSuperAdmin) {
		if user.CompanyID == nil || newUserURIReq.CompanyID != *user.CompanyID {
			log.Errorf("User %d has no company", user.ID)
			c.String(http.StatusForbidden, apierrors.ErrForbidden)
			return
		}
	}

	// Открываем транзакцию
	tx, err := postgres.DB.Beginx()
	if err != nil {
		log.Errorf("Error opening transaction: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	defer tx.Rollback()

	// Проверям, что юзера с такой почтой и таким телефоном нет - иначе 409.
	emailPhoneExtIDIsExists, err := companyqueries.CheckUserEmailPhoneExtIDExists(tx, newUserReq.Email, newUserReq.Phone, newUserReq.ExtID)
	if err != nil {
		log.Errorf("Error checking user email, phone and ext_id: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	if emailPhoneExtIDIsExists {
		log.Errorf("User with this email or phone or ext_id already exists")
		c.String(http.StatusConflict, apierrors.ErrUserEmailPhoneExists)
		return
	}

	if newUserReq.ManagerID != nil {
		_, err := userqueries.GetUserByID(*newUserReq.ManagerID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Errorf("User with this manager_id not found")
				c.String(http.StatusNotFound, "User with this manager_id not found")
				return
			}
			log.Errorf("Error getting user by id: %w", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}
	}

	// # Если передана должность, то проверяем есть ли она в компании
	if newUserReq.PositionID != nil {
		positions, err := companyqueries.GetCompanyPositions(newUserURIReq.CompanyID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Errorf("Company %d has no positions", newUserURIReq.CompanyID)
				c.String(http.StatusNotFound, "Position with this id not found")
				return
			}
			log.Errorf("Error getting position by id: %w", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}
		isFound := false
		for _, position := range positions {
			if position.ID == *newUserReq.PositionID {
				isFound = true
				break
			}
		}
		if !isFound {
			log.Errorf("Position with this id not found")
			c.String(http.StatusNotFound, apierrors.ErrPositionNotFound)
			return
		}
	}

	// # Генерирование токена изменения пароля
	changePasswordTokenLength := viper.GetViper().GetInt("notification.tokenLength")
	changePasswordToken := noauthfunctions.CreateRandomString(changePasswordTokenLength)

	// Сохраняем данные пользователя в БД
	newUserID, err := companyqueries.InsertUser(tx, &newUserReq, newUserURIReq.CompanyID, changePasswordToken)
	if err != nil {
		log.Errorf("Error inserting user: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Обновляем роли пользователя в БД
	if len(newUserReq.RoleIDs) > 0 {
		dbRoles, err := companyqueries.GetRoles(tx, newUserReq.RoleIDs, user.LocalIso)
		if err != nil {
			log.Errorf("Error getting roles: %w", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}
		for i := range newUserReq.RoleIDs {
			isFound := false
			for j := range dbRoles {
				if newUserReq.RoleIDs[i] == dbRoles[j].ID {
					isFound = true
					break
				}
			}
			if !isFound {
				log.Errorf("Role with id %d not found", newUserReq.RoleIDs[i])
				c.String(http.StatusNotFound, apierrors.ErrRoleNotFound)
				return
			}
		}
		// Если все роли на месте, то добавляем их в БД
		err = companyqueries.AddUserRoles(tx, newUserID, newUserReq.RoleIDs)
		if err != nil {
			log.Errorf("Error adding user roles: %w", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}
	}
	// Формируем сообщение, пока связь с БД есть
	title, message, err := composeNewUserNotification(tx, newUserReq, changePasswordToken, user.LocalIso)
	if err != nil {
		log.Errorf("Error composing new user notification: %v", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Error committing transaction: %v", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	bwClient, err := bwclient.NewClient(os.Getenv("BACK_WEB_URL"), 0)
	if err != nil {
		log.Errorf("Error creating back.web client: %v", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Создаем привязки к хозяйствам
	if len(newUserReq.StructureIDs) > 0 {
		statusCode, err := bwClient.SetUserMemberships(newUserURIReq.CompanyID, newUserID, newUserReq.StructureIDs)
		if err != nil {
			log.Errorf("Error setting user memberships: %v", err)
			c.String(statusCode, apierrors.ErrInternalServerError)
			return
		}
	}

	// Создаем привязки к полям
	if len(newUserReq.FieldGUIDs) > 0 {
		statusCode, err := bwClient.SetUserFieldResponsibles(newUserURIReq.CompanyID, newUserID, newUserReq.FieldGUIDs)
		if err != nil {
			log.Errorf("Error setting user field responsibles: %v", err)
			c.String(statusCode, apierrors.ErrInternalServerError)
			return
		}
	}

	// Отправляем сообщение на почту
	err = notificator.SendNotification(newUserReq.Email, title, message)
	if err != nil {
		log.Errorf("Error sending new user notification: %v", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": newUserID})
}

func composeNewUserNotification(tx *sqlx.Tx, newUserReq models.PostNewUserReq, changePasswordToken string, localISO string) (string, string, error) {
	title := ""
	message := ""

	nLocalizations, err := companyqueries.GetLocalizedData(tx, localISO, []int64{viper.GetViper().GetInt64("notification.signUpNotificationTitleId"), viper.GetViper().GetInt64("notification.signUpNotificationMessageId")})
	if err != nil {
		return title, message, err
	}

	for i := range nLocalizations {
		if nLocalizations[i].KeyID == viper.GetViper().GetInt64("notification.signUpNotificationTitleId") {
			title = nLocalizations[i].Value
		} else if nLocalizations[i].KeyID == viper.GetViper().GetInt64("notification.signUpNotificationMessageId") {
			message = nLocalizations[i].Value
		}
	}
	middleName := ""
	if newUserReq.MiddleName != nil {
		middleName = " " + *newUserReq.MiddleName
	}
	message = strings.ReplaceAll(message, "%FULL_NAME%", fmt.Sprintf("%s%s", newUserReq.FirstName, middleName))
	message = strings.ReplaceAll(message, "%SUBMIT_URL%", fmt.Sprintf("%s/auth/change_password?change_password_token=%s&email=%s", os.Getenv("FRONT_URL"), changePasswordToken, newUserReq.Email))
	message = strings.ReplaceAll(message, "%FRONT_URL%", os.Getenv("FRONT_URL"))
	message = strings.ReplaceAll(message, "%GOOGLE_PLAY_APP_LINK%", os.Getenv("GOOGLE_PLAY_APP_LINK"))
	message = strings.ReplaceAll(message, "%CURRENT_YEAR%", fmt.Sprintf("%d", time.Now().Year()))
	message = strings.ReplaceAll(message, "%%TITLE%%", title)

	return title, message, nil
}
