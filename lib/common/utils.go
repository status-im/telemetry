package common

import (
	"database/sql"
	"fmt"
)

func Cleanup(db *sql.DB, table string, before int64) (int64, error) {
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE createdAt < $1", table))
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(before)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
