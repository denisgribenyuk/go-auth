package functions

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	companyqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	noauthqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth/db_queries"
	userqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/user/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/backweb/bwclient"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func PatchUserImage(c *gin.Context) {

	var editUserURIReq models.CompanyUserURIReq

	// Обрабатываем параметры пути
	err := c.ShouldBindUri(&editUserURIReq)
	if err != nil {
		log.Errorf("Error binding uri: %w", err)
		c.String(http.StatusNotFound, apierrors.ErrURLNotFound)
		return
	}

	// Вытаскиваем пользователя, инициировавшего запрос
	user := c.MustGet("user").(*model.User)
	user.LocalIso = strings.ToLower(user.LocalIso)

	// Вытаскиваем пользователя, данные которого нужно изменить
	dbUser, err := userqueries.GetUserByID(editUserURIReq.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("User %d not found in DB", editUserURIReq.UserID)
			c.String(http.StatusNotFound, apierrors.ErrUserNotExists)
			return
		}
		log.Errorf("Error getting user: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Для роли 1 проверяем доступ к компании, для которой создается пользователь. Для роли 3 - можно создавать пользователей для всех компаний
	if !user.HasRole(model.RoleSuperAdmin) {
		if user.CompanyID == nil || editUserURIReq.CompanyID != *user.CompanyID {
			log.Errorf("User %d has no access to company", user.ID)
			c.String(http.StatusForbidden, apierrors.ErrUserHasNoAccessToCompany)
			return
		}

		if dbUser.CompanyID == nil || editUserURIReq.CompanyID != *dbUser.CompanyID {
			log.Errorf("DB user %d has no access to company", dbUser.ID)
			c.String(http.StatusForbidden, apierrors.ErrUserHasNoAccessToCompany)
			return
		}
	}

	if !user.HasRole(model.RoleAdmin) && !user.HasRole(model.RoleSuperAdmin) && user.HasRole(model.RoleUser) {
		if user.ID != dbUser.ID {
			log.Errorf("User %d has no access to user %d", user.ID, dbUser.ID)
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

	// # В зависимости от уровня доступа применяем нужную модель для сериализации
	// if user_is_account_admin or user_is_super_admin:
	if user.HasRole(model.RoleAdmin) || user.HasRole(model.RoleSuperAdmin) {
		// Обрабатываем параметры запроса
		var editUserReq models.EditUserAdminReq

		err = c.ShouldBindJSON(&editUserReq)
		if err != nil {
			log.Errorf("Error binding json: %w", err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if editUserReq.ManagerID != nil {
			manager, err := noauthqueries.GetUserByID(*editUserReq.ManagerID)
			if err != nil {
				if err == sql.ErrNoRows || manager.CompanyID == nil || *manager.CompanyID != editUserURIReq.CompanyID {
					log.Errorf("Manager %d not found in DB", *editUserReq.ManagerID)
					c.String(http.StatusNotFound, apierrors.ErrManagerNotFound)
					return
				}
				log.Errorf("Error getting manager: %w", err)
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			dbUser.ManagerID = editUserReq.ManagerID
		}

		// Проверка, что внешний id не занят
		if editUserReq.ExtID != nil {
			if companyqueries.IsExtIDAlreadyExists(tx, *editUserReq.ExtID) {
				log.Errorf("External User ID %s already exists", *editUserReq.ExtID)
				c.String(http.StatusConflict, apierrors.ErrExternalUserIDExists)
				return
			}
			dbUser.ExtID = editUserReq.ExtID
		}

		// Если передана должность, то проверяем есть ли она в компании

		if editUserReq.PositionID != nil {
			positions, err := companyqueries.GetCompanyPositions(editUserURIReq.CompanyID)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Errorf("Company %d has no positions", editUserURIReq.CompanyID)
					c.String(http.StatusNotFound, "Position with this id not found")
					return
				}
				log.Errorf("Error getting position by id: %w", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}
			isFound := false
			for _, position := range positions {
				if position.ID == *editUserReq.PositionID {
					isFound = true
					break
				}
			}
			if !isFound {
				log.Errorf("Position with this id not found")
				c.String(http.StatusNotFound, apierrors.ErrPositionNotFound)
				return
			}
			dbUser.PositionID = editUserReq.PositionID
		}

		// Обновляем роли пользователя в БД
		if editUserReq.RoleIDs != nil {
			// Can't set role 3 if user is not super admin
			if !user.HasRole(model.RoleSuperAdmin) {
				for _, roleID := range *editUserReq.RoleIDs {
					if roleID == model.RoleSuperAdmin {
						log.Errorf("User %d has no access to set role %d", user.ID, roleID)
						c.String(http.StatusForbidden, apierrors.ErrUserSetRoleForbidden)
						return
					}
				}
			}
			// Сначала удаляем все роли юзера из БД, а затем добавляем эти в БД
			err = companyqueries.DeleteUserRoles(tx, editUserURIReq.UserID)
			if err != nil {
				log.Errorf("Error deleting user roles: %w", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}

			if len(*editUserReq.RoleIDs) > 0 {
				// Проверяем, что все роли существуют в БД
				dbRoles, err := companyqueries.GetRoles(tx, *editUserReq.RoleIDs, user.LocalIso)
				if err != nil {
					log.Errorf("Error getting roles: %w", err)
					c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
					return
				}
				for _, r := range *editUserReq.RoleIDs {
					isFound := false
					for j := range dbRoles {
						if r == dbRoles[j].ID {
							isFound = true
							break
						}
					}
					if !isFound {
						log.Errorf("Role with id %d not found", r)
						c.String(http.StatusNotFound, apierrors.ErrRoleNotFound)
						return
					}
				}

				err = companyqueries.AddUserRoles(tx, editUserURIReq.UserID, *editUserReq.RoleIDs)
				if err != nil {
					log.Errorf("Error adding user roles: %w", err)
					c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
					return
				}
			}
		}

		bwClient, err := bwclient.NewClient(os.Getenv("BACK_WEB_URL"), 0)
		if err != nil {
			log.Errorf("Error creating back.web client: %v", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}

		// Создаем привязки к хозяйствам
		if editUserReq.StructureIDs != nil {
			oldStructures, err := companyqueries.GetUserStructures(tx, editUserURIReq.UserID)
			if err != nil {
				log.Errorf("Error getting user structures: %v", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}

			statusCode, err := bwClient.SetUserMemberships(editUserURIReq.CompanyID, editUserURIReq.UserID, *editUserReq.StructureIDs)
			if err != nil || statusCode != http.StatusOK {
				// Если не получается, пытаемся откатить на старые
				// Return old memberships
				oldStructureIDs := make([]int64, 0)
				for _, s := range oldStructures {
					oldStructureIDs = append(oldStructureIDs, s.StructureID)
				}

				bwClient.SetUserMemberships(editUserURIReq.CompanyID, editUserURIReq.UserID, oldStructureIDs)
				log.Errorf("Error setting user memberships: %v", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}
		}

		// Создаем привязки к полям
		if editUserReq.FieldGUIDs != nil {
			oldFields, err := companyqueries.GetUserFields(tx, editUserURIReq.UserID)
			if err != nil {
				log.Errorf("Error getting user fields: %v", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}
			statusCode, err := bwClient.SetUserFieldResponsibles(editUserURIReq.CompanyID, editUserURIReq.UserID, *editUserReq.FieldGUIDs)
			if err != nil || statusCode != http.StatusOK {
				// Если не получается, пытаемся откатить на старые
				// Return old memberships
				oldFieldGUIDs := make([]string, 0)
				for _, f := range oldFields {
					oldFieldGUIDs = append(oldFieldGUIDs, f.FieldGUID.String())
				}

				bwClient.SetUserFieldResponsibles(editUserURIReq.CompanyID, editUserURIReq.UserID, oldFieldGUIDs)
				log.Errorf("Error setting user fields: %v", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}
		}

		dbUser.FirstName = editUserReq.FirstName
		dbUser.LastName = editUserReq.LastName
		if editUserReq.MiddleName != nil {
			dbUser.MiddleName = editUserReq.MiddleName
		}

		// Проверка, что номер телефона не занят
		if editUserReq.Phone != nil {
			phoneUsers, err := userqueries.GetUserWithPhone(tx, *editUserReq.Phone)
			if err == nil {
				for i := range phoneUsers {
					if phoneUsers[i].ID != editUserURIReq.UserID {
						log.Errorf("Phone number %d already exists", *editUserReq.Phone)
						c.String(http.StatusConflict, apierrors.ErrPhoneAlreadyExists)
						return
					}
				}
			}
			dbUser.Phone = *editUserReq.Phone
		}
		if editUserReq.ActiveFlag {
			dbUser.ActiveFlag = 1
		} else {
			dbUser.ActiveFlag = 0
		}
	} else {
		var editUserReq models.EditUserReq

		err = c.ShouldBindJSON(&editUserReq)
		if err != nil {
			log.Errorf("Error binding json: %w", err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		//     # Для роли 2 - нельзя редактировать данные неактивного пользователя
		if dbUser.ActiveFlag == 0 {
			log.Errorf("User %d is not active", editUserURIReq.UserID)
			c.String(http.StatusForbidden, apierrors.ErrUserNotActive)
			return
		}

		dbUser.FirstName = editUserReq.FirstName
		dbUser.LastName = editUserReq.LastName
		if editUserReq.MiddleName != nil {
			dbUser.MiddleName = editUserReq.MiddleName
		}

		// Проверка, что номер телефона не занят
		if editUserReq.Phone != nil {
			phoneUsers, err := userqueries.GetUserWithPhone(tx, *editUserReq.Phone)
			if err == nil {
				for i := range phoneUsers {
					if phoneUsers[i].ID != editUserURIReq.UserID {
						log.Errorf("Phone number %d already exists", *editUserReq.Phone)
						c.String(http.StatusConflict, apierrors.ErrPhoneAlreadyExists)
						return
					}
				}
			}
			dbUser.Phone = *editUserReq.Phone
		}
	}

	// # Обновление данных пользователя
	err = companyqueries.UpdateUser(tx, dbUser)
	if err != nil {
		log.Errorf("Error updating user: %v", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Закрыавем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Errorf("Error committing transaction: %v", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	c.String(http.StatusOK, "OK")
}
