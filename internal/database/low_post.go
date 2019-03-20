package database

import (
	"database/sql"
	"escapade/internal/models"

	"time"
	//
	_ "github.com/lib/pq"
)

func postPath(tx *sql.Tx, id int) (path string, err error) {
	query := `select path from
							Post where id = $1
						 `
	row := tx.QueryRow(query, id)

	err = row.Scan(&path)
	return
}

// postCreate create post
func (db *DataBase) postCreate(tx *sql.Tx, post models.Post, thread models.Thread, t time.Time) (createdPost models.Post, err error) {

	var path string
	if post.Parent == 0 {
		path = ""
	} else {
		if path, err = postPath(tx, post.Parent); err != nil {
			return
		}
	}
	query := `INSERT INTO Post(author, created, forum, message, thread, parent, path) VALUES
						 	($1, $2, $3, $4, $5, $6, $7) 
						 RETURNING id, author, created, forum, message, thread, parent, path;
						 `
	row := tx.QueryRow(query, post.Author, t,
		thread.Forum, post.Message, thread.ID, post.Parent, path)

	createdPost = models.Post{}
	if err = row.Scan(&createdPost.ID, &createdPost.Author, &createdPost.Created,
		&createdPost.Forum, &createdPost.Message, &createdPost.Thread, &createdPost.Parent,
		&createdPost.Path); err != nil {
		return
	}
	return
}

func (db *DataBase) postsGetFlat(tx *sql.Tx, thread models.Thread, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundPosts []models.Post, err error) {

	query := `select id, author, created, forum, message, thread, parent from
							Post where thread = $1 and lower(forum) like lower($2)`

	if tb {
		if desc {
			query += ` and created <= $3`
			query += ` order by created, id desc`
		} else {
			query += ` and created >= $3`
			query += ` order by created, id`
		}
		if lb {
			query += ` Limit $4`
		}
	} else if lb {
		if desc {
			query += ` order by created, id desc`
		} else {
			query += ` order by created, id`
		}
		query += ` Limit $3`
	}

	var rows *sql.Rows

	if tb {
		if lb {
			rows, err = tx.Query(query, thread.ID, thread.Forum, t, limit)
		} else {
			rows, err = tx.Query(query, thread.ID, thread.Forum, t)
		}
	} else if lb {
		rows, err = tx.Query(query, thread.ID, thread.Forum, limit)
	} else {
		rows, err = tx.Query(query, thread.ID, thread.Forum)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {

		post := models.Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.Forum, &post.Message, &post.Thread, &post.Parent); err != nil {
			break
		}
		foundPosts = append(foundPosts, post)
	}
	return
}
