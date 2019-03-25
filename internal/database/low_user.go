package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createUser(tx *sql.Tx, user *models.User) (createdUser models.User, err error) {

	query := `INSERT INTO UserForum(fullname, nickname, email, about) VALUES
						 	($1, E'` + user.Nickname + `', $2, $3) 
						 `
	queryAddUserReturning(&query)
	row := tx.QueryRow(query, user.Fullname, user.Email, user.About)

	createdUser = models.User{}
	if err = row.Scan(&createdUser.ID, &createdUser.Fullname, &createdUser.Nickname,
		&createdUser.Email, &createdUser.About); err != nil {
		return
	}
	return
}

func (db *DataBase) updateUser(tx *sql.Tx, user models.User) (updated models.User, err error) {

	query := `	UPDATE UserForum set fullname = $1, email = $2, about = $3
								where nickname = $4
						`
	queryAddUserReturning(&query)
	row := tx.QueryRow(query, user.Fullname, user.Email, user.About, user.Nickname)

	updated = models.User{}
	if err = row.Scan(&updated.ID, &updated.Fullname, &updated.Nickname,
		&updated.Email, &updated.About); err != nil {
		return
	}
	return
}

func (db *DataBase) userCheckID(tx *sql.Tx, oldNickname string) (newNickname string, err error) {
	var thatUser models.User
	if thatUser, err = db.findUserByName(tx, oldNickname); err != nil {
		err = re.ErrorUserNotExist()
		return
	}
	newNickname = thatUser.Nickname
	return
}

// confirmUnique confirm that user.Email and user.Name
// dont use by another Player
func (db DataBase) userConfirmUnique(tx *sql.Tx, user *models.User) (users *[]models.User, err error) {

	var foundUsers *[]models.User
	users = &[]models.User{}

	foundUsers, err = db.isOnlyEmailUnique(tx, user.Email, user.Nickname)

	if err != nil {
		return
	}

	if foundUsers != nil && len(*foundUsers) > 0 {
		*users = append(*users, *foundUsers...)
		foundUsers = nil
	}

	foundUsers, err = db.isNicknameUnique(tx, user.Nickname)

	if err != nil {
		return
	}

	if foundUsers != nil && len(*foundUsers) > 0 {
		*users = append(*users, *foundUsers...)
		foundUsers = nil
	}

	return
}

func (db DataBase) findUser(tx *sql.Tx, queryAdd string, arg string) (foundUser models.User, err error) {

	query := `SELECT fullname, nickname, email, about 
	FROM UserForum ` + queryAdd

	row := tx.QueryRow(query, arg)

	foundUser = models.User{}
	if err = row.Scan(&foundUser.Fullname, &foundUser.Nickname,
		&foundUser.Email, &foundUser.About); err != nil {
		return
	}
	return
}

func (db DataBase) findUserByName(tx *sql.Tx, taken string) (foundUser models.User, err error) {

	query := `where lower(nickname) like lower($1)`
	foundUser, err = db.findUser(tx, query, taken)
	return
}

func (db DataBase) findUserByEmail(tx *sql.Tx, taken string) (foundUser models.User, err error) {

	query := `where lower(email) like lower($1)`
	foundUser, err = db.findUser(tx, query, taken)
	return
}

/*
query:
	select fullname, nickname, email, about
		FROM UserForum as uf
		where
		lower(nickname) > lower(E'UO5l6xBUiVhHJ.bill') and
		(
		nickname in
		(
			select author
				from Post
				where
				lower(uf.nickname) like lower(author) ESCAPE '' and
				lower(forum) like lower($1)
		) or
		nickname in
		(
			select author
				from Thread
				where
				lower(uf.nickname) like lower(author) ESCAPE '' and
				lower(forum) like lower($1)
		)
	)
		order by lower(nickname)
	  Limit 4


*/

func (db *DataBase) usersGet(tx *sql.Tx, slug string,
	qc QueryGetConditions) (foundUsers []models.User, err error) {

	pq := &postQuery{
		sortASC:     ` order by lower(nickname) `,
		sortDESC:    ` order by lower(nickname) desc `,
		compareASC:  ` and lower(nickname) > lower('` + qc.nv + `')`,
		compareDESC: ` and lower(nickname) < lower('` + qc.nv + `')`,
	}

	query := `
	select fullname, nickname, email, about 
		FROM UserForum as uf 
		where (
			nickname in 
		(
			select author
				from Post
				where 
				lower(uf.nickname) like lower(author) and
				lower(forum) like lower($1)
		) or
		nickname in 
		(
			select author
				from Thread
				where 
				lower(uf.nickname) like lower(author) and
				lower(forum) like lower($1)
		)
		)
	`
	queryAddConditions(&query, qc, pq)

	fmt.Println("query:" + query)
	var rows *sql.Rows

	if rows, err = tx.Query(query, slug); err != nil {
		return
	}
	defer rows.Close()

	foundUsers = []models.User{}
	for rows.Next() {

		user := models.User{}
		if err = rows.Scan(&user.Fullname, &user.Nickname,
			&user.Email, &user.About); err != nil {
			break
		}
		foundUsers = append(foundUsers, user)
	}

	return
}

func (db DataBase) findUsers(tx *sql.Tx, queryAdd string, taken ...string) (foundUsers *[]models.User, err error) {

	query := `SELECT fullname, nickname, email, about 
	FROM UserForum ` + queryAdd

	var rows *sql.Rows

	if len(taken) == 1 {
		rows, err = tx.Query(query, taken[0])
	} else {
		rows, err = tx.Query(query, taken[0], taken[1])
	}

	if err != nil {
		return
	}
	defer rows.Close()

	foundUsers = &[]models.User{}
	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.Fullname, &user.Nickname,
			&user.Email, &user.About); err != nil {
			break
		}

		*foundUsers = append(*foundUsers, user)
	}
	return
}

// isNicknameUnique checks if there are Players with
// this('taken') nickname and returns corresponding error if yes. A
func (db DataBase) isNicknameUnique(tx *sql.Tx, taken string) (foundUsers *[]models.User, err error) {

	query := `where lower(nickname) like lower($1)`
	foundUsers, err = db.findUsers(tx, query, taken)
	return
}

// isEmailUnique checks if there are Players with
// this('taken') email and returns corresponding error if yes. B\A
func (db DataBase) isOnlyEmailUnique(tx *sql.Tx, email string, nickname string) (foundUsers *[]models.User, err error) {

	query := `where lower(email) like lower($1) and lower(nickname) not like lower($2)`
	foundUsers, err = db.findUsers(tx, query, email, nickname)
	return
}

// query add returning
func queryAddUserReturning(query *string) {
	*query += ` RETURNING id, fullname, nickname, email, about `
}

// scan row to model Vote
func userScan(row *sql.Row) (foundUser models.User, err error) {
	foundUser = models.User{}
	err = row.Scan(&foundUser.Fullname, &foundUser.Nickname,
		&foundUser.Email, &foundUser.About)
	return
}
