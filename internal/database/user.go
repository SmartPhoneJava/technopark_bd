package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) CreateUser(user *models.User) (foundUsers *[]models.User, createdUser models.User, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if foundUsers, err = db.userConfirmUnique(tx, user); err != nil || len(*foundUsers) > 0 {
		return
	}

	if createdUser, err = db.createUser(tx, user); err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (db *DataBase) GetUser(name string) (foundUser models.User, err error) {
	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if foundUser, err = db.findUserByName(tx, name); err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (db *DataBase) UpdateUser(user models.User) (foundUser models.User, err error) {
	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if foundUser, err = db.findUserByEmail(tx, user.Email); err != nil {
		if err != sql.ErrNoRows {
			return
		}
	}

	if foundUser.Nickname != "" && foundUser.Nickname != user.Nickname {
		err = re.ErrorEmailIstaken()
		return
	}

	if foundUser, err = db.updateUser(tx, user); err != nil {
		return
	}
	err = tx.Commit()
	return
}
