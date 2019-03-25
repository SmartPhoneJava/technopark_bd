package rerrors

import "errors"

func ErrorInvalidPath() error {
	return errors.New("Path invalid")
}

func ErrorInvalidID() error {
	return errors.New("id invalid")
}

func ErrorPostConflict() error {
	return errors.New("Parent post was created in another thread")
}
