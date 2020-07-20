package errors

import "fmt"

type ErrorType int

func Type(ce error) ErrorType {
	if e, ok := ce.(CustomError); ok {
		return e.Code
	}
	return Unknown
}

func (e ErrorType) NewWithMessage(msg string) CustomError {
	return CustomError{
		Code:    e,
		Message: msg,
	}
}

func (e ErrorType) NewWithMessageF(tmpl string, args ...interface{}) CustomError {
	msg := fmt.Sprintf(tmpl, args...)
	return CustomError{
		Code:    e,
		Message: msg,
	}
}

func (e ErrorType) Wrap(err error, msg string) CustomError {
	return CustomError{
		RootCause: err,
		Code:      e,
		Message:   msg,
	}
}

func (e ErrorType) WrapF(err error, tmpl string, args ...interface{}) CustomError {
	msg := fmt.Sprintf(tmpl, args...)
	return e.Wrap(err, msg)
}
