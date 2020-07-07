package errors

const (
	ProjectNotFound ErrorType = -(1200 + iota)
	ProjectQueryError
	ProjectPermissionQueryError
	ProjectPermissionNotFound
	ProjectPermissionExisted
	ProjectPermissionCreatingError
	ProjectCreatingError
	ProjectCannotUpdate
	ProjectCannotDelete
)
