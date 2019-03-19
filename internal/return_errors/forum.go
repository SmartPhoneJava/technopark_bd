package rerrors

import "errors"

func ErrorForumSlugIsTaken() error {
	return errors.New("Slug is taken")
}

func ErrorForumUserNotExist() error {
	return errors.New("User not exist")
}
