package rerrors

import "errors"

func ErrorVoteNotExist() error {
	return errors.New("Thread not exist")
}
