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

	if returnThread, err = db.threadConfirmUnique(tx, thread); err != nil {
		return
	}

	if thread.Author, err = db.userCheckID(tx, thread.Author); err != nil {
		return
	}

	if thread.Forum, err = db.forumCheckID(tx, thread.Forum); err != nil {
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

	if returnThreads, err = db.threadsGet(tx, slug, limit, existLimit, t, existTime, desc); err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *DataBase) GetThread(slug string) (returnThread models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}

	err = tx.Commit()
	return
}
