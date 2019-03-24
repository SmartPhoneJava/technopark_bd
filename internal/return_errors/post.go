package rerrors

import "errors"

func ErrorInvalidPath() error {
	return errors.New("Path invalid")
}
