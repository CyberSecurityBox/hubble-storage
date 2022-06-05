package postgresql

import (
	"database/sql"
)

func (psqlconn *Psqlconn) SelectSourcePodName(sourcePodName string) ([]string, error) {
	var selectQuery = "select metadata " +
		"from flows where source_pod_name=$1"
	var err error
	var result []string
	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = psqlconn.Db.Prepare(selectQuery)
	if err != nil {
		psqlconn.Log.WithError(err).Error("Can't select db.Prepare")
		return result, err
	}
	rows, err = stmt.Query(sourcePodName)
	if err != nil {
		psqlconn.Log.WithError(err).Error("Can't select db.Query")
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var metadata string

		err = rows.Scan(&metadata)
		if err != nil {
			psqlconn.Log.WithError(err).Error("Can't select rows.Scan")
		}

		result = append(result, metadata)
	}

	return result, nil
}