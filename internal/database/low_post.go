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
func (db *DataBase) postCreate(tx *sql.Tx, post models.Post, thread models.Thread, t time.Time) (createdPost models.Post, err error) {

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


1.5
1.6.7
1.8.8

*/

//
func queryAddConditions(queryInit string, qc QueryGetConditions,
	sortASC string, sortDESC string, compareIDASC string, compareIDDESC string) (query string) {
	query = queryInit
	if qc.desc {
		if qc.mn {
			query += compareIDDESC
		}
		query += sortDESC
	} else {
		if qc.mn {
			query += compareIDASC
		}
		query += sortASC
	}
	if qc.ln {
		query += ` Limit ` + strconv.Itoa(qc.lv)
	}
	return //23 -> 19
}

/*
var path string
		if path, _, err = getPathAndLvl(tx, qc.mv); err != nil {
			return
		} // ` and id > ` + strconv.Itoa(qc.mv) + `` //
		query += compareID //` and path > '` + path + `'`
*/

func (db *DataBase) postsGetFlat(tx *sql.Tx, thread models.Thread, slug string,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {
	asc := " order by created, id "
	desc := " order by created desc, id desc "
	compareIDASC := `and id > ` + strconv.Itoa(qc.mv)
	compareIDDESC := `and id < ` + strconv.Itoa(qc.mv)
	foundPosts, err = db.postsGet(tx, thread, slug, qc, asc, desc, compareIDASC, compareIDDESC)
	return
}

func (db *DataBase) postsGetTree(tx *sql.Tx, thread models.Thread, slug string,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {
	asc := " order by string_to_array(path, '.')::int[], created "
	desc := " order by	string_to_array(path, '.')::int[] desc, created desc "

	var path string

	if qc.mn {
		if path, _, err = getPathAndLvl(tx, qc.mv); err != nil {
			return
		}
	}

	fmt.Println("after:" + path)
	compareIDASC := ` and path > '` + path + `'`
	compareIDDESC := ` and path < '` + path + `'`
	foundPosts, err = db.postsGet(tx, thread, slug, qc, asc, desc, compareIDASC, compareIDDESC)
	return
}

func (db *DataBase) postsGetParentTree(tx *sql.Tx, thread models.Thread, slug string,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {
	asc := ` order by string_to_array(path, '.')::int[], created `
	desc := ` order by split_part(path, '.', 2) desc, string_to_array(path, '.')::int[], created `

	var path string

	if qc.mn {
		if path, _, err = getPathAndLvl(tx, qc.mv); err != nil {
			return
		}
	}

	fmt.Println("afterpostsGetParentTree:" + path)
	compareIDASC := ` and path > '` + path + `'`
	compareIDDESC := ` and path < '` + path + `' and split_part(path, '.', 2) < split_part('` + path + `', '.', 2)`

	if qc.lv, err = parentTreeGetLimit(tx, thread, qc,
		asc, desc, compareIDASC, compareIDDESC); err != nil {
		fmt.Println("parentTreeGetLimit err")
		return
	}

	foundPosts, err = db.postsGet(tx, thread, slug, qc, asc, desc, compareIDASC, compareIDDESC)
	return
}

func parentTreeGetLimit(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, sortASC string, sortDESC string,
	compareIDASC string, compareIDDESC string) (realLimit int, err error) {

	queryInside := `select * from Post where 1 = 1 
	`

	limitNeed := qc.ln
	qc.ln = false
	queryInside = queryAddConditions(queryInside, qc, sortASC, sortDESC, compareIDASC, compareIDDESC)

	qc.ln = limitNeed

	query := `select COUNT(*), split_part(path, '.', 2) from ( ` + queryInside + ` ) as A
					GROUP BY split_part(path, '.', 2), forum, thread
					HAVING thread = $1 and lower(forum) like lower($2)
					`

	if qc.desc {
		query += "order by split_part(path, '.', 2) desc"
	} else {
		query += "order by split_part(path, '.', 2)"
	}

	query += ` Limit $3`

	var rows *sql.Rows

	fmt.Println("prepare query:", query)

	if rows, err = tx.Query(query, thread.ID, thread.Forum, qc.lv); err != nil {
		return
	}
	defer rows.Close()

	fmt.Println("was:", qc.lv)

	realLimit = 0
	for rows.Next() {
		var count int
		var str string
		if err = rows.Scan(&count, &str); err != nil {
			break
		}
		fmt.Println("str:" + str)
		realLimit += count
	}
	if err != nil {
		return
	}

	fmt.Println("now:", realLimit)
	return
}

func (db *DataBase) postsGet(tx *sql.Tx, thread models.Thread, slug string,
	qc QueryGetConditions, sortASC string, sortDESC string, compareIDASC string, compareIDDESC string) (foundPosts []models.Post, err error) {

	query := `select id, author, created, forum, message, thread, parent, path, level from
							Post where thread = $1 and lower(forum) like lower($2)`

	query = queryAddConditions(query, qc, sortASC, sortDESC, compareIDASC, compareIDDESC)

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

// 370. lets do 200 - 161 done
