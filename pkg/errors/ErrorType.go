package errors

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

func (e ErrorType) Wrap(err error, msg string) CustomError {
	return CustomError{
		RootCause: err,
		Code:      e,
		Message:   msg,
	}
}
