package errors

const (
	DatasetNotFound ErrorType = -(1400 + iota)
	DatasetQueryError
	DatasetCannotCreate
	DatasetCannotDelete
)
