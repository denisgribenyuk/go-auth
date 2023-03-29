package dbqueries

import (
	"github.com/jmoiron/sqlx"
	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

// GetUserByImageID returns user by image ID
func GetUserByImageID(tx *sqlx.Tx, imageID int64) ([]*model.User, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.first_name,
			u.last_name,
			u.middle_name,
			u.local_iso,
			u.company_id,
			u.phone,
			u.image_id,
			u.ext_id,
			u.manager_id,
			u.position_id
		FROM
			users AS u
		WHERE
			u.image_id = $1
	`
	var res []*model.User
	err := tx.Select(&res, query, imageID)
	return res, err
}

// InsertUserImage inserts user image into DB and returns its ID
func InsertUserImage(tx *sqlx.Tx, imagePath string) (int64, error) {
	query := `
		INSERT
		INTO
			images (image_path)
		VALUES
			($1)
		RETURNING id
	`
	var id int64
	err := postgres.DB.Get(&id, query, imagePath)
	return id, err
}

// UpdateUserImage updates user image in DB
func UpdateUserImage(tx *sqlx.Tx, userID, imageID int64) error {
	query := `
		UPDATE
			users
		SET
			image_id = $1
		WHERE id = $2
	`
	_, err := tx.Exec(query, imageID, userID)
	return err
}

// DeleteUserImage deletes user image from DB
func DeleteUserImage(tx *sqlx.Tx, imageID int64) error {
	query := `
		DELETE
		FROM
			images
		WHERE id = $1
	`
	_, err := tx.Exec(query, imageID)
	return err
}
