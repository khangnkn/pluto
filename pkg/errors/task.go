package errors

const (
	TaskCannotCreate ErrorType = -(1600 + iota)
	TaskCannotGet
	TaskCannotDelete
	TaskNotFound
	TaskDetailCannotGet
	TaskDetailCannotUpdate
	TaskDetailCannotDelete
)
