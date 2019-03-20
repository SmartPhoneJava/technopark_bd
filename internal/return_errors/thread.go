package rerrors

import "errors"

func ErrorThreadNotExist() error {
	return errors.New("Thread not exist")
}

func ErrorThreadConflict() error {
	return errors.New("Thread is taken")
}
