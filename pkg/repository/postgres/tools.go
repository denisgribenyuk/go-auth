package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetCount(tx *sqlx.Tx, query string) (int, error) {
	preparedQuery := fmt.Sprintf("SELECT count(*) AS count FROM (%s) cte", query)
	var count int
	if tx != nil {
		err := tx.Get(&count, preparedQuery)
		if err != nil {
			return 0, err
		}
		return count, nil
	} else if DB != nil {
		err := DB.Get(&count, preparedQuery)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, fmt.Errorf("db connection is not initialized")
}
