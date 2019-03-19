package rerrors

import "errors"

func ErrorForumSlugIsTaken() error {
	return errors.New("Slug is taken")
}
