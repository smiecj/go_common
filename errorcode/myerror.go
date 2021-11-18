package errorcode

import (
	"fmt"
)

// 自定义error类
type myError struct {
	code string
	msg  string
}

// to string
func (err myError) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", err.code, err.msg)
}

// 判断error 是否完全一致，可用于 errors.Is 判断
func (err myError) Is(compareError error) bool {
	transformErr, ok := compareError.(myError)
	if !ok {
		return false
	}
	return err.code == transformErr.code
}

// build error with msg
func BuildErrorWithMsg(code, msg string) error {
	return myError{code: code, msg: msg}
}

// build error only with code
func BuildError(code string) error {
	return myError{code: code, msg: code}
}
