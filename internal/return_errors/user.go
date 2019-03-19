package rerrors

import "errors"

func ErrorUserNotExist() error {
	return errors.New("User not exist")
}
