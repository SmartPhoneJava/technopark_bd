package rerrors

import "errors"

func ErrorThreadNotExist() error {
	return errors.New("Thread not exist")
}
