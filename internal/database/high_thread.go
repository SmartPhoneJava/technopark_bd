package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"
	"time"

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

func (db *DataBase) GetThreads(slug string, limit int, existLimit bool, t time.Time, existTime bool, desc bool) (returnThreads []models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if _, err = db.findForumBySlug(tx, slug); err != nil {
		err = re.ErrorForumNotExist()
		return
	}

	fmt.Println("GetThreads got:", t.String())
	if returnThreads, err = db.threadsGet(tx, slug, limit, existLimit, t, existTime, desc); err != nil {
		return
	}

	err = tx.Commit()
	return
}
