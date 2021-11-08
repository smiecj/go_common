package errorcode

import "fmt"

type myError struct {
	code string
	msg  string
}

func (err myError) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", err.code, err.msg)
}

func BuildError(code, msg string) error {
	return myError{code: code, msg: msg}
}
