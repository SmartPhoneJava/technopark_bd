package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) voteCreate(tx *sql.Tx, vote models.Vote, thread models.Thread) (createdVote models.Vote, err error) {

	query := `INSERT INTO Vote(author, voice, thread) VALUES
						 	($1, $2, $3) 
						 RETURNING author, voice, thread;
						 `
	row := tx.QueryRow(query, vote.Author, vote.Voice, thread.ID)

	createdVote = models.Vote{}
	if err = row.Scan(&createdVote.Author, &createdVote.Voice,
		&createdVote.Thread); err != nil {
		return
	}
	return
}

func (db DataBase) voteFindByThreadAndAuthor(tx *sql.Tx, thread int, author string) (foundVote models.Vote, err error) {

	query := `SELECT voice, thread, author, isEdited FROM Vote where thread = $1 and author = $2`

	row := tx.QueryRow(query, thread, author)

	foundVote = models.Vote{}
	if err = row.Scan(&foundVote.Voice, &foundVote.Thread, &foundVote.Author, &foundVote.IsEdited); err != nil {
		return
	}
	return
}

func (db DataBase) voteUpdate(tx *sql.Tx, vote models.Vote, thread models.Thread) (updatedVote models.Vote, err error) {

	query := `	UPDATE Vote set voice = $1--, isEdited = true
		where author = $2 and thread = $3 and isEdited = false
		RETURNING author, voice, thread;
	`

	row := tx.QueryRow(query, vote.Voice, vote.Author, thread.ID)

	updatedVote = models.Vote{}
	if err = row.Scan(&updatedVote.Author, &updatedVote.Voice,
		&updatedVote.Thread); err != nil {
		return
	}
	return
}

func (db *DataBase) voteThread(tx *sql.Tx, voice int, thread models.Thread) (updated models.Thread, err error) {

	query := `	UPDATE Thread set votes = votes + $1
								where id = $2
								RETURNING id, slug, author, created, forum, message, title, votes
						 `

	row := tx.QueryRow(query, voice, thread.ID)

	updated = models.Thread{}
	if err = row.Scan(&updated.ID, &updated.Slug, &updated.Author,
		&updated.Created, &updated.Forum, &updated.Message,
		&updated.Title, &updated.Votes); err != nil {
		return
	}
	return
}

/*
Author   string `json:"nickname" db:"author"`
Voice    int    `json:"voice" db:"voice"`
Thread   int    `json:"-" db:"thread"`
IsEdited bool   `json:"-" db:"isEdited"`
*/
