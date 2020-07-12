package errors

type CustomError struct {
	RootCause error
	Code      ErrorType
	Message   string
}

func (ce CustomError) Error() string {
	if ce.RootCause != nil {
		return ce.Message + " >> " + ce.RootCause.Error()
	}
	return ce.Message
}

func (ce CustomError) BareError() string {
	return ce.Message
}
