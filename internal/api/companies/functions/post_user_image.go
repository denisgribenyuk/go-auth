package functions

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path"

	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	dbqueries "gitlab.assistagro.com/back/back.auth.go/internal/api/companies/db_queries"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

// PostUserImage create/change user image
func PostUserImage(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	var userImageURI models.CompanyUserURIReq
	var userImageQuery models.PostUserImageQueryReq

	{
		// Обрабатываем параметры пути
		err := c.ShouldBindUri(&userImageURI)
		if err != nil {
			log.Errorf("Error binding uri: %w", err)
			c.String(http.StatusNotFound, apierrors.ErrURLNotFound)
			return
		}

		// Обрабатываем параметры запроса
		err = c.ShouldBindJSON(&userImageQuery)
		if err != nil {
			log.Errorf("Error binding query: %w", err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	// Поехали.
	if *userImageQuery.Base64 == "" {
		c.Render(
			http.StatusOK, render.Data{
				ContentType: "text/plain",
				Data:        []byte(*userImageQuery.Base64),
			})
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

	userImage := []byte("")
	if *userImageQuery.Base64 != "" {
		// Декодируем изображение
		userImage, err = base64.StdEncoding.DecodeString(*userImageQuery.Base64)
		if err != nil {
			log.Errorf("Error decoding image: %w", err)
			c.String(http.StatusBadRequest, apierrors.ErrBase64FormatError)
			return
		}
	}
	// Make thumbnail for image with longer side at 200px

	img, _, err := image.Decode(bytes.NewReader(userImage))
	if err != nil {
		log.Errorf("Error decoding image: %w", err)
		c.String(http.StatusBadRequest, apierrors.ErrBase64FormatError)
		return
	}

	// Генерим маленькую иконку с максимальной стороной 200px
	thumbImage := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	// Декодируем основное изображение в jpeg
	imgBuf := new(bytes.Buffer)
	err = jpeg.Encode(imgBuf, img, nil)
	if err != nil {
		log.Errorf("Error encoding image: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Декодируем иконку в jpeg
	thumbBuf := new(bytes.Buffer)
	err = jpeg.Encode(thumbBuf, thumbImage, nil)
	if err != nil {
		log.Errorf("Error encoding image: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Генерим уникальные имена файлов
	imgName := fmt.Sprintf("%s.jpeg", uuid.New())
	imgThumbName := fmt.Sprintf("thumb_%s", imgName)

	// Вытаскиваем переменные среды окружения для бакета и каталога в бакете, в который класть файлы
	awsBucket := os.Getenv("S3_AUTH_BUCKET_NAME")
	photosPholder := viper.GetString("s3.photos_folder")

	err = repository.S3.PutObject(awsBucket, path.Join(photosPholder, imgName), imgBuf.Bytes())
	if err != nil {
		log.Errorf("Error putting image to s3: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	repository.S3.PutObject(awsBucket, path.Join(photosPholder, imgThumbName), thumbBuf.Bytes())
	if err != nil {
		log.Errorf("Error putting image to s3: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// # Были случаи, когда у нескольких пользователей установлена фотка с одинаковым id.
	// # Это нужно обработать, поэтому смотрим количество пользователей с этой фоткой.

	count := 0
	if dbUser.ImageID != nil {
		countUserWitImageID, err := dbqueries.GetUserByImageID(tx, *dbUser.ImageID)
		if err != nil {
			log.Errorf("Error getting user by image id: %w", err)
			c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
			return
		}
		count = len(countUserWitImageID)
	}

	// Добавляем изображение в БД
	newImageID, err := dbqueries.InsertUserImage(tx, imgName)
	if err != nil {
		log.Errorf("Error inserting user image: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// Обновляем информацию об изображении у пользователя
	err = dbqueries.UpdateUserImage(tx, dbUser.ID, newImageID)
	if err != nil {
		log.Errorf("Error updating user image: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	// #  Если старое фото задано только одному пользователю, то удаляем запись с фото из БД
	if count <= 1 {
		if dbUser.ImageID != nil {
			err = dbqueries.DeleteUserImage(tx, *dbUser.ImageID)
			if err != nil {
				log.Errorf("Error deleting user image: %w", err)
				c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
				return
			}
		}
		// Плюем на ошибки, не удалось удалить фото - не страшно
		repository.S3.DeleteObject(awsBucket, dbUser.ImagePath)
	}
	err = tx.Commit()

	if err != nil {
		log.Errorf("Error commit transaction: %w", err)
		c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, imgName)
}
