package base

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (err *Error) Error() string { // 指针接收者
	return err.Message
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func NewErrorf(code int, format string, a ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, a...)}
}
