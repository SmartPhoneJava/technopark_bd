package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

// CreateThread handle thread creation
func (db *DataBase) CreatePost(posts []models.Post) (createdPosts []models.Post, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	createdPosts = []models.Post{}
	for _, post := range posts {
		// if returnForum, err = db.postConfirmUnique(tx, forum); err != nil {
		// 	return
		// }

		if post.Author, err = db.userCheckID(tx, post.Author); err != nil {
			return
		}

		if post.Forum, err = db.forumCheckID(tx, post.Forum); err != nil {
			return
		}

		if post.Thread, err = db.threadCheckID(tx, post.Thread); err != nil {
			return
		}

		if post, err = db.postCreate(tx, post); err != nil {
			return
		}

		createdPosts = append(createdPosts, post)
	}
	err = tx.Commit()
	return
}
