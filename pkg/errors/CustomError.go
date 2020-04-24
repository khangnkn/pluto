package errors

type CustomError struct {
	RootCause error
	Code      ErrorType
	Message   string
}

func (ce CustomError) Error() string {
	return ce.Message + " >> " + ce.RootCause.Error()
}

func (ce CustomError) BareError() string {
	return ce.Message
}
