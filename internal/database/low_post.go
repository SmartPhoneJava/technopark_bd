package database

import (
	"database/sql"
	"escapade/internal/models"

	"time"
	//
	_ "github.com/lib/pq"
)

/*
id SERIAL PRIMARY KEY,
        author varchar(120) not null,
        forum varchar(120),
        message varchar(1600) not null,
        created    TIMESTAMPTZ,
        isEdited boolean default false,
        thread int ,
        parent int
*/

// postCreate create post
func (db *DataBase) postCreate(tx *sql.Tx, post models.Post, thread models.Thread) (createdPost models.Post, err error) {

	query := `INSERT INTO Post(author, created, forum, message, thread) VALUES
						 	($1, $2, $3, $4, $5) 
						 RETURNING id, author, created, forum, message, thread;
						 `
	row := tx.QueryRow(query, post.Author, time.Now(),
		thread.Forum, post.Message, thread.ID)

	createdPost = models.Post{}
	if err = row.Scan(&createdPost.ID, &createdPost.Author, &createdPost.Created,
		&createdPost.Forum, &createdPost.Message, &createdPost.Thread); err != nil {
		return
	}
	return
}
