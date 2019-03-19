package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"time"

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

/*
SELECT 	a.FieldWidth, a.FieldHeight,
					a.MinsTotal, a.MinsFound,
					a.Finished, a.Exploded
	 FROM Player as p
		JOIN
			(
				SELECT player_id,
					FieldWidth, FieldHeight,
					MinsTotal, MinsFound,
					Finished, Exploded
					FROM Game Order by id
			) as a
			ON p.id = a.player_id and p.name like $1
			OFFSET $2 Limit $3

*/

// getThreads get threads
func (db *DataBase) threadsGetWithLimit(tx *sql.Tx, slug string, limit int) (foundThreads []models.Thread, err error) {

	query := `select id, slug, author, created, forum, message, title from
							Thread where forum like $1 Limit $2;
						 `

	var rows *sql.Rows

	if rows, err = tx.Query(query, slug, limit); err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {
		thread := models.Thread{}
		if err = rows.Scan(&thread.ID, &thread.Slug,
			&thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title); err != nil {
			break
		}

		foundThreads = append(foundThreads, thread)
	}
	return
}

func (db *DataBase) threadsGetWithLimitAndTime(tx *sql.Tx, slug string, limit int, t time.Time) (foundThreads []models.Thread, err error) {

	query := `select id, slug, author, created, forum, message, title from
							Thread where forum like $1 and created like $2 Limit $3;
						 `

	var rows *sql.Rows

	if rows, err = tx.Query(query, slug, t, limit); err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {
		thread := models.Thread{}
		if err = rows.Scan(&thread.ID, &thread.Slug,
			&thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title); err != nil {
			break
		}

		foundThreads = append(foundThreads, thread)
	}
	return
}

func (db *DataBase) threadsGetWithTime(tx *sql.Tx, slug string, t time.Time) (foundThreads []models.Thread, err error) {

	query := `select id, slug, author, created, forum, message, title from
							Thread where forum like $1 and created like $2 ;
						 `

	var rows *sql.Rows

	if rows, err = tx.Query(query, slug, t); err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {
		thread := models.Thread{}
		if err = rows.Scan(&thread.ID, &thread.Slug,
			&thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title); err != nil {
			break
		}

		foundThreads = append(foundThreads, thread)
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

// threadCheckForum checks, is thread's forum exists
func (db *DataBase) threadCheckForum(tx *sql.Tx, thread *models.Thread) (err error) {
	var thatForum models.Forum
	if thatForum, err = db.findForumBySlug(tx, thread.Forum); err != nil {
		err = re.ErrorForumNotExist()
		return
	}
	thread.Forum = thatForum.Slug
	return
}
