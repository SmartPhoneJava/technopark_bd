package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"

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
		fmt.Println("sorry")
		return
	}

	var thatUser models.User
	if thatUser, err = db.findUserByName(tx, forum.User); err != nil {
		err = re.ErrorForumUserNotExist()
		return
	}
	forum.User = thatUser.Nickname

	if returnForum, err = db.createForum(tx, forum); err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (db *DataBase) GetForum(slug string) (returnForum models.Forum, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnForum, err = db.findForumBySlug(tx, slug); err != nil {
		err = re.ErrorForumUserNotExist()
		return
	}

	err = tx.Commit()
	return
}
