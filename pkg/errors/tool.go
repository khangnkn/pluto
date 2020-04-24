package errors

const (
	ToolNotFound ErrorType = -(1000 + iota)
	ToolNoRecord
	ToolQueryError
)
