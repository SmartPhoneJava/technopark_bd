package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) CreateForum(forum *models.Forum) (returnForum models.Forum, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnForum, err = db.forumConfirmUnique(tx, forum); err != nil {
		return
	}

	if returnForum, err = db.createForum(tx, forum); err != nil {
		return
	}
	err = tx.Commit()
	return
}
