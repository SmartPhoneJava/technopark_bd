package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

// createThread create thread
func (db *DataBase) threadCreate(tx *sql.Tx, thread *models.Thread) (createdThread models.Thread, err error) {

	query := `INSERT INTO Thread(slug, author, created, forum, message, title) VALUES
						 	($1, $2, $3, $4, $5, $6) 
						 RETURNING id, slug, author, created, forum, message, title;
						 `
	row := tx.QueryRow(query, thread.Slug, thread.Author, thread.Created,
		thread.Forum, thread.Message, thread.Title)

	createdThread = models.Thread{}
	if err = row.Scan(&createdThread.ID, &createdThread.Slug,
		&createdThread.Author, &createdThread.Created, &createdThread.Forum,
		&createdThread.Message, &createdThread.Title); err != nil {
		return
	}
	return
}

// checkUser checks, is thread's author exists
func (db *DataBase) threadCheckUser(tx *sql.Tx, thread *models.Thread) (err error) {
	var thatUser models.User
	if thatUser, err = db.findUserByName(tx, thread.Author); err != nil {
		err = re.ErrorUserNotExist()
		return
	}
	thread.Author = thatUser.Nickname
	return
}

// checkUser checks, is thread's forum exists
func (db *DataBase) threadCheckForum(tx *sql.Tx, thread *models.Thread) (err error) {
	var thatForum models.Forum
	if thatForum, err = db.findForumBySlug(tx, thread.Forum); err != nil {
		err = re.ErrorForumNotExist()
		return
	}
	thread.Forum = thatForum.Slug
	return
}
