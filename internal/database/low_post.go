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

type postQuery struct {
	sortASC     string
	sortDESC    string
	compareASC  string
	compareDESC string
}

func getPathAndLvl(tx *sql.Tx, id int) (path string, lvl int, err error) {
	query := `select path, level from Post where id = $1
						 `

	if err = tx.QueryRow(query, id).Scan(&path, &lvl); err != nil {
		return
	}
	return
}

func updatePath(path *string, id int) {
	*path = *path + "." + strconv.Itoa(id)
}

// postCreate create post
func (db *DataBase) postCreate(tx *sql.Tx, post models.Post, thread models.Thread,
	t time.Time) (createdPost models.Post, err error) {

	var (
		path string
		lvl  int
	)
	if post.Parent == 0 {
		path = "0"
		lvl = 0
	} else {
		if path, lvl, err = getPathAndLvl(tx, post.Parent); err != nil {
			return
		}
		if path == "" {
			err = re.ErrorInvalidPath()
		}
	}
	lvl++

	query := `INSERT INTO Post(author, created, forum, message, thread, parent, path, level) VALUES
						 	($1, $2, $3, $4, $5, $6, $7, $8) 
						 RETURNING id, author, created, forum, message, thread, parent, path, level;
						 `
	row := tx.QueryRow(query, post.Author, t,
		thread.Forum, post.Message, thread.ID, post.Parent, path, lvl)

	createdPost = models.Post{}
	if err = row.Scan(&createdPost.ID, &createdPost.Author, &createdPost.Created,
		&createdPost.Forum, &createdPost.Message, &createdPost.Thread, &createdPost.Parent,
		&createdPost.Path, &createdPost.Level); err != nil {
		return
	}

	query = `UPDATE Post set path=$1 where id=$2
						`
	updatePath(&path, createdPost.ID)
	_, err = tx.Exec(query, path, createdPost.ID)

	return
}

/*
// for threads

func queryAddConditions(queryInit string, qc QueryGetConditions, sortASC string, sortDESC string) (query string) {
	query = queryInit
	if qc.tn {
		if qc.desc {
			query += ` and created <= $3`
			query += sortDESC
		} else {
			query += ` and created >= $3`
			query += sortASC
		}
		if qc.ln {
			query += ` Limit $4`
		}
	} else if qc.ln {
		if qc.desc {
			query += sortDESC
		} else {
			query += sortASC
		}
		query += ` Limit $3`
	}
	return //23 -> 19
}
*/

func queryAddConditions(query *string, qc QueryGetConditions, pq *postQuery) {
	queryAddMinID(query, qc, pq.compareASC, pq.compareDESC)
	queryAddSort(query, qc, pq.sortASC, pq.sortDESC)
	queryAddLimit(query, qc)
} // 22

func queryAddSort(query *string, qc QueryGetConditions, sortASC string, sortDESC string) {
	if qc.desc {
		*query += sortDESC
	} else {
		*query += sortASC
	}
}

func queryAddMinID(query *string, qc QueryGetConditions, compareIDASC string, compareIDDESC string) {
	if qc.mn {
		if qc.desc {
			*query += compareIDDESC
		} else {
			*query += compareIDASC
		}
	}
}

func queryAddLimit(query *string, qc QueryGetConditions) {
	if qc.ln {
		*query += ` Limit ` + strconv.Itoa(qc.lv)
	}
	return
}

func (db *DataBase) postsGetFlat(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	pq := &postQuery{
		sortASC:     ` order by created, id `,
		sortDESC:    ` order by created desc, id desc `,
		compareASC:  `and id > ` + strconv.Itoa(qc.mv),
		compareDESC: `and id < ` + strconv.Itoa(qc.mv),
	}
	foundPosts, err = db.postsGet(tx, thread, qc, pq)
	return
}

func (db *DataBase) postsGetTree(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	var path string

	if qc.mn {
		if path, _, err = getPathAndLvl(tx, qc.mv); err != nil {
			return
		}
	}
	pq := &postQuery{
		sortASC: ` order by string_to_array(path, '.')::int[], created `,
		sortDESC: ` order by	string_to_array(path, '.')::int[] desc, created desc `,
		compareASC:  ` and path > '` + path + `'`,
		compareDESC: ` and path < '` + path + `'`,
	}
	foundPosts, err = db.postsGet(tx, thread, qc, pq)
	return
}

func (db *DataBase) postsGetParentTree(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	var path string

	if qc.mn {
		if path, _, err = getPathAndLvl(tx, qc.mv); err != nil {
			return
		}
	}
	pq := &postQuery{
		sortASC:    ` order by string_to_array(path, '.')::int[], created `,
		sortDESC:   ` order by split_part(path, '.', 2) desc, string_to_array(path, '.')::int[], created `,
		compareASC: ` and path > '` + path + `'`,
		compareDESC: ` and path < '` + path + `
			' and split_part(path, '.', 2) < split_part('` + path + `', '.', 2)`,
	}

	foundPosts, err = parentTreeGet(tx, thread, qc, pq)
	return
}

func (db *DataBase) postsGet(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, pq *postQuery) (foundPosts []models.Post, err error) {

	query := `
	select id, author, created, forum, message,
		thread, parent, path, level 
		 from Post 
		 where thread = $1 and 
			 lower(forum) like lower($2)
	`
	queryAddConditions(&query, qc, pq)

	fmt.Println("query:" + query)
	var rows *sql.Rows

	if rows, err = tx.Query(query, thread.ID, thread.Forum); err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {

		post := models.Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.Forum, &post.Message, &post.Thread, &post.Parent,
			&post.Path, &post.Level); err != nil {
			break
		}
		foundPosts = append(foundPosts, post)
	}

	return
}

func parentTreeGet(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, pq *postQuery) (foundPosts []models.Post, err error) {

	queryInside := `
		select split_part(path, '.', 2) as p
			from Post 
				where thread = $1 and lower(forum) like lower($2)
	`
	queryAddMinID(&queryInside, qc, pq.compareASC, pq.compareDESC)

	groupBy := `
		select A.p from 
		( 
			` + queryInside + ` 
		) as A
		GROUP BY A.p
	`
	queryAddSort(&groupBy, qc, "order by A.p", "order by A.p desc")
	queryAddLimit(&groupBy, qc)

	query := `
	select id, author, created, forum,
		message, thread, parent, path, level 
			from Post 
			where thread = $1 and lower(forum) like lower($2)
				 	and split_part(path, '.', 2) = ANY 
				 	(` + groupBy + `)
	`
	queryAddSort(&query, qc, pq.sortASC, pq.sortDESC)

	var rows *sql.Rows

	if rows, err = tx.Query(query, thread.ID, thread.Forum); err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {

		post := models.Post{}
		if err = rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.Forum, &post.Message, &post.Thread, &post.Parent,
			&post.Path, &post.Level); err != nil {
			break
		}
		foundPosts = append(foundPosts, post)
	}

	return
}

// 280
