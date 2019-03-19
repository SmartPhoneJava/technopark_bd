package rerrors

import "errors"

func ErrorForumSlugInvalid() error {
	return errors.New("Slug invalid")
}

func ErrorForumSlugIsTaken() error {
	return errors.New("Slug is taken")
}

func ErrorForumNotExist() error {
	return errors.New("Forum not exist")
}
