package errors

const (
	TaskCannotCreate ErrorType = -(1600 + iota)
	TaskCannotGet
	TaskCannotDelete
	TaskDetailCannotGet
	TaskDetailCannotUpdate
	TaskDetailCannotDelete
)
