package errors

const (
	ImageNotFound ErrorType = -(1500 + iota)
	ImageQueryError
	ImageTooManyRequest
)
