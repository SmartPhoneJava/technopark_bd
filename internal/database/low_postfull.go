package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) postfullGet(tx *sql.Tx, existRelated bool, related string, id int) (fullpost models.Postfull, err error) {

	query := `
	select id, author, created, forum,
	message, thread, parent, path, level
		 from Post 
		 where id = $1
	`
	fullpost.Post = &models.Post{}
	if err = tx.QueryRow(query, id).Scan(&fullpost.Post.ID, &fullpost.Post.Author, &fullpost.Post.Created,
		&fullpost.Post.Forum, &fullpost.Post.Message, &fullpost.Post.Thread, &fullpost.Post.Parent,
		&fullpost.Post.Path, &fullpost.Post.Level); err != nil {
		return
	}

	return
}
