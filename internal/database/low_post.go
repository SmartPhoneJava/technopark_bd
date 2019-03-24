package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"
	"strconv"

	"time"
	//
	_ "github.com/lib/pq"
)

func postPath(tx *sql.Tx, id int) (path string, err error) {
	query := `select path from Post where id = $1
						 `
	row := tx.QueryRow(query, id)

	if err = row.Scan(&path); err != nil {
		fmt.Println("ooops")
		return
	}

	updatePath(&path, id)
	return
}

func updatePath(path *string, id int) {
	*path = *path + "." + strconv.Itoa(id)
}

// postCreate create post
func (db *DataBase) postCreate(tx *sql.Tx, post models.Post, thread models.Thread, t time.Time) (createdPost models.Post, err error) {

	var path string
	if post.Parent == 0 {
		path = "0"
	} else {
		if path, err = postPath(tx, post.Parent); err != nil {
			return
		}
		if path == "" {
			err = re.ErrorInvalidPath()
			fmt.Println("no path")
		}
	}

	// concat_ws('.', $7, id::text))

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
		fmt.Println("no create")
		return
	}

	query = `UPDATE Post set path=$1 where id=$2
						`
	updatePath(&path, createdPost.ID)
	_, err = tx.Exec(query, path, createdPost.ID)

	return
}

func (db *DataBase) postsGetFlat(tx *sql.Tx, thread models.Thread, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundPosts []models.Post, err error) {

	query := `select id, author, created, forum, message, thread, parent from
							Post where thread = $1 and lower(forum) like lower($2)`

	if tb {
		if desc {
			query += ` and created <= $3`
			query += ` order by created desc, id desc`
		} else {
			query += ` and created >= $3`
			query += ` order by created, id`
		}
		if lb {
			query += ` Limit $4`
		}
	} else if lb {
		if desc {
			query += ` order by created desc, id desc`
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
		fmt.Println("sorry100")
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

func (db *DataBase) postsGetTree(tx *sql.Tx, thread models.Thread, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundPosts []models.Post, err error) {

	query := `select id, author, created, forum, message, thread, parent, path from
							Post where thread = $1 and lower(forum) like lower($2)`

	if tb {
		if desc {
			query += ` and created <= $3`
			query += ` order by
												  string_to_array(path, '.')::int[] desc, created desc;
							`
		} else {
			query += ` and created >= $3`
			query += ` order by
													string_to_array(path, '.')::int[], created`
		}
		if lb {
			query += ` Limit $4`
		}
	} else if lb {
		if desc {
			query += ` order by 
													string_to_array(path, '.')::int[] desc, created desc`
		} else {
			query += ` order by 
													string_to_array(path, '.')::int[], created`
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
		fmt.Println("sorry163")
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {

		post := models.Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.Forum, &post.Message, &post.Thread, &post.Parent,
			&post.Path); err != nil {
			fmt.Println("wrong path")
			break
		}
		foundPosts = append(foundPosts, post)
	}
	return
}

func (db *DataBase) postsGetParentTree(tx *sql.Tx, thread models.Thread, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundPosts []models.Post, err error) {

	query := `select id, author, created, forum, message, thread, parent, path from
							Post where thread = $1 and lower(forum) like lower($2)`

	if tb {
		if desc {
			query += ` and created <= $3`
			query += ` order by
			split_part(path, '.', 2) desc, string_to_array(path, '.')::int[], created;
							`
		} else {
			query += ` and created >= $3`
			query += ` order by
													string_to_array(path, '.')::int[], created`
		}
		if lb {
			query += ` Limit $4`
		}
	} else if lb {
		if desc {
			query += ` order by 
				split_part(path, '.', 2) desc, string_to_array(path, '.')::int[], created
				`
		} else {
			query += ` order by 
													string_to_array(path, '.')::int[], created`
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
		fmt.Println("sorry163")
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {

		post := models.Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.Forum, &post.Message, &post.Thread, &post.Parent,
			&post.Path); err != nil {
			fmt.Println("wrong path")
			break
		}
		foundPosts = append(foundPosts, post)
	}
	return
}

type QueryParameters struct {
	query  string
	thread int
	forum  int
}

type QueryConditions struct {
	tv   time.Time // time value
	tn   bool      // time need
	lv   int       // limit value
	ln   bool      // limit need
	desc bool      // desc need
}

func queryAddCondition(query *string, yes bool, add string) {
	if yes {
		*query = *query + add
	}
}

/*
func queryAddConditions(query *string, qc QueryConditions, sort string) {
	if qc.tn {
		queryAddCondition(query, qc.desc, " and created <= $3")
		queryAddCondition(query, !qc.desc, " and created >= $3")
		queryAddCondition(query, qc.desc,, sort + " desc")
		queryAddCondition(query, !qc.desc,, sort)
		queryAddCondition(query, qc.ln, " Limit $4")
	} else {

		queryAddCondition(query, qc.desc,, sort + " desc")
		queryAddCondition(query, !qc.desc,, sort)
		queryAddCondition(query, qc.ln, " Limit $3")
	}
}
2019-03-24T14:23:03.905Z
2019-03-24T14:23:03.953Z

2019-03-24T14:25:07.693Z"
2019-03-24T14:25:07.736Z
*/

func (db *DataBase) postsGetUsual(tx *sql.Tx, thread models.Thread, slug string, limit int, lb bool,
	t time.Time, tb bool, desc bool) (foundPosts []models.Post, err error) {

	query := `select id, author, created, forum, message, thread, parent, path from
							Post where thread = $1 and lower(forum) like lower($2)`

	if tb {
		if desc {
			query += ` and created <= $3`
			query += ` order by
												  string_to_array(path, '.')::int[], created desc;
							`
		} else {
			query += ` and created >= $3`
			query += ` order by
													string_to_array(path, '.')::int[], created`
		}
		if lb {
			query += ` Limit $4`
		}
	} else if lb {
		if desc {
			query += ` order by 
													string_to_array(path, '.')::int[], created desc`
		} else {
			query += ` order by 
													string_to_array(path, '.')::int[], created`
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
		fmt.Println("sorry163")
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {

		post := models.Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.Forum, &post.Message, &post.Thread, &post.Parent,
			&post.Path); err != nil {
			fmt.Println("wrong path")
			break
		}
		foundPosts = append(foundPosts, post)
	}
	return
}
