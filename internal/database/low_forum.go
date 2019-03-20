package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createForum(tx *sql.Tx, forum *models.Forum) (createdForum models.Forum, err error) {

	query := `INSERT INTO Forum(slug, title, user_nickname ) VALUES
						 	($1, $2, $3) 
						 RETURNING slug, title, user_nickname ;
						 `
	row := tx.QueryRow(query, forum.Slug, forum.Title, forum.User)

	createdForum = models.Forum{}
	if err = row.Scan(&createdForum.Slug, &createdForum.Title, &createdForum.User); err != nil {
		return
	}
	fmt.Println("i can show slug:", createdForum.Slug)
	return
}

func (db *DataBase) forumCheckID(tx *sql.Tx, oldNickname string) (newNickname string, err error) {
	var thatForum models.Forum
	if thatForum, err = db.findForumBySlug(tx, oldNickname); err != nil {
		err = re.ErrorForumNotExist()
		return
	}
	newNickname = thatForum.Slug
	return
}

func (db DataBase) forumConfirmUnique(tx *sql.Tx, forum *models.Forum) (foundForum models.Forum, err error) {

	if foundForum, err = db.findForumBySlug(tx, forum.Slug); err != nil && err != sql.ErrNoRows {
		return
	}
	err = nil

	if foundForum.Slug != "" {
		err = re.ErrorForumSlugIsTaken()
		return
	}

	return
}

func (db DataBase) findForum(tx *sql.Tx, queryAdd string, arg string) (foundForum models.Forum, err error) {

	query := `SELECT slug, title, user_nickname 
		FROM Forum ` + queryAdd

	row := tx.QueryRow(query, arg)

	foundForum = models.Forum{}
	if err = row.Scan(&foundForum.Slug, &foundForum.Title, &foundForum.User); err != nil {
		return
	}
	return
}

func (db DataBase) findForumBySlug(tx *sql.Tx, taken string) (foundForum models.Forum, err error) {

	query := `where lower(slug) like lower($1)`
	foundForum, err = db.findForum(tx, query, taken)
	return
}
