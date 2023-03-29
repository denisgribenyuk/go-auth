package functions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetUserImage(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(*model.User)
	var userImageURI models.CompanyUserURIReq
	var userImageQuery models.GetUserImageQueryReq

	// Обрабатываем параметры пути
	err := c.ShouldBindUri(&userImageURI)
	if err != nil {
		log.Errorf("Error binding uri: %w", err)
		c.String(http.StatusNotFound, apierrors.ErrURLNotFound)
		return
	}

	// Обрабатываем параметры запроса
	err = c.ShouldBindQuery(&userImageQuery)
	if err != nil {
		log.Errorf("Error binding query: %w", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	tx, err := postgres.DB.Beginx()
	if err != nil {
		log.Errorf("Error begin transaction: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}
	defer tx.Rollback()

	// Получаем пользователя из БД
	dbUser, err := dbqueries.GetUserImage(tx, userImageURI.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, apierrors.ErrUserImageNotFound)
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Проверяем, есть ли у пользователя доступ к компании пользователя
	if !user.HasRole(model.RoleSuperAdmin) {
		if dbUser.CompanyID == nil || *dbUser.CompanyID != userImageURI.CompanyID {
			c.String(http.StatusForbidden, apierrors.ErrUserNotInCompany)
			return
		}
	}

	if dbUser.ImageID == nil {
		c.String(http.StatusNotFound, apierrors.ErrUserImageNotFound)
		return
	}
	res := models.UserPhoto{
		Name: dbUser.ImagePath,
	}

	// # Если запрашивают Thumbnail
	if *userImageQuery.IsThumbnail == 1 {
		res.Name = "thumb_" + res.Name
	}

	awsBucket := os.Getenv("S3_AUTH_BUCKET_NAME")
	photosPholder := viper.GetString("s3.photos_folder")
	imgObj := path.Join(photosPholder, res.Name)
	fmt.Println("Key: ", imgObj)
	imageBytes, err := repository.S3.GetObject(awsBucket, imgObj)
	if err != nil {
		log.Errorf("Image is not available on S3. %v", err)
		c.String(http.StatusNotFound, fmt.Sprintf("Image is not available on S3. %v", err))
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Error commit transaction: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	res.Image = base64.StdEncoding.EncodeToString(imageBytes)

	c.JSON(http.StatusOK, res)
}
