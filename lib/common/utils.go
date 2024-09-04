package common

import (
	"database/sql"
	"fmt"
)

func Cleanup(db *sql.DB, table string, before int64) (int64, error) {
	result, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE createdAt < $1", table), before)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
