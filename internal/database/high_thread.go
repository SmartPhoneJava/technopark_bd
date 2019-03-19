package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

// CreateThread handle thread creation
func (db *DataBase) CreateThread(thread *models.Thread) (returnThread models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	// if returnForum, err = db.threadConfirmUnique(tx, forum); err != nil {
	// 	return
	// }

	if err = db.threadCheckUser(tx, thread); err != nil {
		return
	}

	if err = db.threadCheckForum(tx, thread); err != nil {
		return
	}

	if returnThread, err = db.threadCreate(tx, thread); err != nil {
		return
	}
	err = tx.Commit()
	return
}
