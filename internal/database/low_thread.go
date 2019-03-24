package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"strconv"
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

func (db *DataBase) threadsGet(tx *sql.Tx, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundThreads []models.Thread, err error) {

	query := `select id, slug, author, created, forum, message, title from
							Thread where lower(forum) like lower($1)`

	if tb {
		if desc {
			query += ` and created <= $2`
			query += ` order by created desc`
		} else {
			query += ` and created >= $2`
			query += ` order by created`
		}
		if lb {
			query += ` Limit $3`
		}
	} else if lb {
		if desc {
			query += ` order by created desc`
		} else {
			query += ` order by created`
		}
		query += ` Limit $2`
	}

	var rows *sql.Rows

	if tb {
		if lb {
			rows, err = tx.Query(query, slug, t, limit)
		} else {
			rows, err = tx.Query(query, slug, t)
		}
	} else if lb {
		rows, err = tx.Query(query, slug, limit)
	} else {
		rows, err = tx.Query(query, slug)
	}

	if err != nil {
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

func (db DataBase) threadConfirmUnique(tx *sql.Tx, thread *models.Thread) (foundThread models.Thread, err error) {
	// if foundThread, err = db.threadFindByTitle(tx, thread.Title); err != sql.ErrNoRows {
	// 	err = re.ErrorThreadConflict()
	// 	return
	// }
	if thread.Slug != "" {
		if foundThread, err = db.threadFindBySlug(tx, thread.Slug); err != sql.ErrNoRows {
			err = re.ErrorThreadConflict()
			return
		}
	}
	err = nil
	return
}

func (db DataBase) threadFindByTitle(tx *sql.Tx, title string) (foundThread models.Thread, err error) {

	query := `SELECT id, slug, author, created, forum, message, title from
	Thread where title like $1`

	row := tx.QueryRow(query, title)

	foundThread = models.Thread{}
	if err = row.Scan(&foundThread.ID, &foundThread.Slug,
		&foundThread.Author, &foundThread.Created, &foundThread.Forum,
		&foundThread.Message, &foundThread.Title); err != nil {
		return
	}
	return
}

func (db DataBase) threadFindBySlug(tx *sql.Tx, slug string) (foundThread models.Thread, err error) {

	query := `SELECT id, slug, author, created, forum, message, title from
	Thread where lower(slug) like lower($1)`

	row := tx.QueryRow(query, slug)

	foundThread = models.Thread{}
	if err = row.Scan(&foundThread.ID, &foundThread.Slug,
		&foundThread.Author, &foundThread.Created, &foundThread.Forum,
		&foundThread.Message, &foundThread.Title); err != nil {
		return
	}
	return
}

func (db DataBase) threadFindByID(tx *sql.Tx, arg int) (foundThread models.Thread, err error) {

	query := `SELECT id, slug, author, created, forum, message, title from
	Thread where id like $1`

	row := tx.QueryRow(query, arg)

	foundThread = models.Thread{}
	if err = row.Scan(&foundThread.ID, &foundThread.Slug,
		&foundThread.Author, &foundThread.Created, &foundThread.Forum,
		&foundThread.Message, &foundThread.Title); err != nil {
		return
	}
	return
}

func (db *DataBase) threadCheckID(tx *sql.Tx, oldID int) (newID int, err error) {
	var thatThread models.Thread
	if thatThread, err = db.threadFindByID(tx, oldID); err != nil {
		err = re.ErrorThreadNotExist()
		return
	}
	newID = thatThread.ID
	return
}

func (db DataBase) threadFindByIDorSlug(tx *sql.Tx, arg string) (foundThread models.Thread, err error) {

	var (
		id  int
		row *sql.Row
	)
	query := `SELECT id, slug, author, created, forum, message, title from Thread`
	if id, err = strconv.Atoi(arg); err != nil {
		query += ` where lower(slug) like lower($1)`
		row = tx.QueryRow(query, arg)
		err = nil
	} else {
		query += ` where id = $1`
		row = tx.QueryRow(query, id)
	}

	foundThread = models.Thread{}
	if err = row.Scan(&foundThread.ID, &foundThread.Slug,
		&foundThread.Author, &foundThread.Created, &foundThread.Forum,
		&foundThread.Message, &foundThread.Title); err != nil {
		//err = re.ErrorThreadNotExist()
		return
	}
	return
}
