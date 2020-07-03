package errors

const (
	ProjectNotFound ErrorType = -(1200 + iota)
	ProjectQueryError
	ProjectPermissionQueryError
	ProjectPermissionNotFound
	ProjectPermissionExisted
	ProjectCreatingError
	ProjectPermissionCreatingError
)
