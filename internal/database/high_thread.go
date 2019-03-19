package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
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

func (db *DataBase) GetThreads(slug string, limit int, existLimit bool, t time.Time, existTime bool) (returnThreads []models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if _, err = db.findForumBySlug(tx, slug); err != nil {
		err = re.ErrorForumNotExist()
		return
	}

	if existLimit && existTime {
		if returnThreads, err = db.threadsGetWithLimitAndTime(tx, slug, limit, t); err != nil {
			return
		}
	} else if existLimit {
		if returnThreads, err = db.threadsGetWithLimit(tx, slug, limit); err != nil {
			return
		}
	} else if existTime {
		if returnThreads, err = db.threadsGetWithTime(tx, slug, t); err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}
