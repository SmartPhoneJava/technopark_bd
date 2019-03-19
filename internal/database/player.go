package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"

	//
	_ "github.com/lib/pq"
)

// В будущем добавить, чтобы отдельно была проверка на
// на корректность, отдельно на sql  инъекции
// func ValidatePrivateUI(user *models.UserPrivateInfo) (err error) {

// 	if !models.ValidatePassword(user.Password) {
// 		err = re.ErrorInvalidPassword()
// 		return
// 	}

// 	if !models.ValidatePlayerName(user.Name) && !models.ValidateEmail(user.Email) {
// 		err = re.ErrorInvalidNameOrEmail()
// 		return
// 	}

// 	return
// }

// GetPlayerIDbyName get player's id by his hame
func (db *DataBase) GetPlayerIDbyName(username string) (id int, err error) {
	sqlStatement := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlStatement, username)

	err = row.Scan(&id)
	return
}

// GetPlayerNamebyID get player's name by his id
func (db *DataBase) GetPlayerNamebyID(id int) (username string, err error) {
	sqlStatement := `SELECT name FROM Player WHERE id = $1`
	row := db.Db.QueryRow(sqlStatement, id)

	err = row.Scan(&username)
	return
}

// GetNameByEmail get player's name by his email
func (db DataBase) GetNameByEmail(email string) (name string, err error) {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := db.Db.QueryRow(sqlStatement, email)

	if err = row.Scan(&name); err != nil {
		return
	}
	return
}

// GetNameByEmail get player's name by his email
func (db DataBase) GetPasswordEmailByName(name string) (email string, password string, err error) {
	sqlStatement := "SELECT email, password " +
		"FROM Player where name like $1"

	row := db.Db.QueryRow(sqlStatement, name)

	if err = row.Scan(&email, &password); err != nil {
		return
	}
	return
}

// confirmUnique confirm that user.Email and user.Name
// dont use by another Player
func (db DataBase) ConfirmUnique(tx *sql.Tx, user *models.User) (users *[]models.User, err error) {

	var foundUsers *[]models.User
	users = &[]models.User{}

	if foundUsers, err = db.isOnlyEmailUnique(tx, user.Email, user.Nickname); err != nil || len(*foundUsers) > 0 {
		if foundUsers == nil {
			return
		}
		fmt.Println("foundUsers:", len(*foundUsers))

		*users = append(*users, *foundUsers...)
	}
	//foundUsers = nil
	// if foundUsers, err = db.isFullnameUnique(tx, user.Fullname); err != nil || len(*foundUsers) > 0 {
	// 	fmt.Println("foundUsers:", len(*foundUsers))
	// 	*users = append(*users, *foundUsers...)
	// }
	foundUsers = nil
	if foundUsers, err = db.isNicknameUnique(tx, user.Nickname); err != nil || len(*foundUsers) > 0 {
		fmt.Println("foundUsers:", len(*foundUsers))
		*users = append(*users, *foundUsers...)
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

	query := `where nickname like $1`
	foundUser, err = db.findUser(tx, query, taken)
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
	fmt.Println("foundUsers:", len(*foundUsers))
	return
}

// isFullnameUnique checks if there are Players with
// this('taken') fullname and returns corresponding error if yes

// func (db DataBase) isFullnameUnique(tx *sql.Tx, taken string) (foundUsers *[]models.User, err error) {
// 	query := `where lower(fullname)=($1)`
// 	foundUsers, err = db.findUsers(tx, query, taken)
// 	return
// }

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

func (db DataBase) checkBunch(field string, password string) (id int, err error) {
	// If checkBunchNamePass cant find brunch name-password
	if id, err = db.checkBunchNamePass(field, password); err != nil {
		// and checkBunchEmailPass cant find brunch email-password
		if id, err = db.checkBunchEmailPass(field, password); err != nil {
			err = re.ErrorWrongPassword()
			return // then password wrong
		}
	}
	err = nil
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchNamePass(username string, password string) (id int, err error) {

	sqlStatement := "SELECT id FROM Player where name like $1 and password like $2"
	row := db.Db.QueryRow(sqlStatement, username, password)
	err = row.Scan(&id)
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchEmailPass(email string, password string) (id int, err error) {
	sqlStatement := "SELECT id FROM Player where email like $1 and password like $2"
	row := db.Db.QueryRow(sqlStatement, email, password)
	err = row.Scan(&id)
	return
}

// confirmRightEmail checks that Player with such
// email and name exists
// func (db DataBase) confirmEmailNamePassword(user *models.UserPrivateInfo) error {
// 	sqlStatement := "SELECT 1 FROM Player where name like $1 and password like $2 and email like $3"

// 	row := db.Db.QueryRow(sqlStatement, user.Name, user.Password, user.Email)
// 	var res int
// 	err := row.Scan(&res)
// 	return err
// }

// func (db *DataBase) deletePlayer(user *models.UserPrivateInfo) error {
// 	sqlStatement := `
// 	DELETE FROM Player where name=$1 and password=$2 and email=$3
// 		`
// 	_, err := db.Db.Exec(sqlStatement, user.Name,
// 		user.Password, user.Email)

// 	return err
// }

func (db *DataBase) createPlayer(tx *sql.Tx, user *models.User) (createdUser models.User, err error) {

	query := `INSERT INTO UserForum(fullname, nickname, email, about) VALUES
						 	($1, $2, $3, $4) 
						 RETURNING id, fullname, nickname, email, about;
						 `
	row := tx.QueryRow(query, user.Fullname, user.Nickname, user.Email, user.About)

	createdUser = models.User{}
	if err = row.Scan(&createdUser.ID, &createdUser.Fullname, &createdUser.Nickname,
		&createdUser.Email, &createdUser.About); err != nil {
		return
	}
	return
}

//UpdatePlayerByName gets name of Player from
//relation Session, cause we know that user has session
func (db *DataBase) updateUser(tx *sql.Tx, user models.User) (updated models.User, err error) {
	// var (
	// 	curEmail     string
	// 	curAbout      string
	// 	sqlStatement string
	// 	oldName      string
	// )

	// oldName = curName
	// if curEmail, curPass, err = db.GetPasswordEmailByName(curName); err != nil {
	// 	return
	// }

	// if user.Email != curEmail && user.Email != "" {
	// 	if !models.ValidateEmail(user.Email) {
	// 		return re.ErrorInvalidEmail()
	// 	}
	// 	if err = db.isEmailUnique(user.Email); err != nil {
	// 		return re.ErrorInvalidEmail()
	// 	}
	// 	curEmail = user.Email
	// }

	// if user.Password != curPass && user.Password != "" {
	// 	curPass = user.Password
	// }

	// if user.Name != curName && user.Name != "" {
	// 	if !models.ValidateString(user.Name) {
	// 		return re.ErrorInvalidName()
	// 	}
	// 	if err = db.isNameUnique(user.Name); err != nil {
	// 		return re.ErrorInvalidName()
	// 	}
	// 	curName = user.Name
	// }

	query := `	UPDATE UserForum set fullname = $1, email = $2, about = $3
								where nickname = $4
								RETURNING id, fullname, nickname, email, about;
						`
	row := tx.QueryRow(query, user.Fullname, user.Email, user.About, user.Nickname)

	updated = models.User{}
	if err = row.Scan(&updated.ID, &updated.Fullname, &updated.Nickname,
		&updated.Email, &updated.About); err != nil {
		return
	}
	return
}
