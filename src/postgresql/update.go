package postgresql

import (
	"database/sql"
	"strings"
)

func (psqlconn *Psqlconn) UpdateFlowCount(checksum string) error {
	var updateStmt = "update flows " +
		"set flow_count=$1 where checksum=$2"
	var selectQuery = "select flow_count " +
		"from flows where checksum=$1"
	var err error
	var flowCount int
	var prepare *sql.Stmt
	var row *sql.Row

	// Select flow_count with checksum
	prepare, err = psqlconn.Db.Prepare(selectQuery)
	if err != nil {
		psqlconn.Log.WithError(err).Error("Can't select db.Prepare")
		return err
	}
	row = prepare.QueryRow(checksum)
	err = row.Scan(&flowCount)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			psqlconn.Log.WithError(err).Error("Can't select db.Scan")
		}
		return err
	}

	// Update flow_count
	flowCount++
	_, err = psqlconn.Db.Exec(updateStmt, flowCount, checksum)
	if err != nil {
		psqlconn.Log.WithError(err).Error("Can't insert db.Exec")
		return err
	}

	return nil
}
