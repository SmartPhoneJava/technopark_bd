package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"

	//
	_ "github.com/lib/pq"
)

// CreateThread handle thread creation
func (db *DataBase) CreatePost(posts []models.Post, slug string) (createdPosts []models.Post, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	createdPosts = []models.Post{}

	var thatThread models.Thread

	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		fmt.Println("forum noooooooooo exists:")
		return
	}

	fmt.Println("forum exists:", thatThread.ID, thatThread.Slug)

	for _, post := range posts {
		// if returnForum, err = db.postConfirmUnique(tx, forum); err != nil {
		// 	return
		// }

		if post.Author, err = db.userCheckID(tx, post.Author); err != nil {
			return
		}

		if post, err = db.postCreate(tx, post, thatThread); err != nil {
			return
		}

		createdPosts = append(createdPosts, post)
	}
	err = tx.Commit()
	return
}
