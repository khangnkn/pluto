package errors

const (
	TaskCannotCreate ErrorType = -(1600 + iota)
	TaskCannotGet
	TaskDetailCannotGet
	TaskDetailCannotUpdate
	TaskDetailCannotDelete
)
