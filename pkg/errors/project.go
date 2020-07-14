package errors

const (
	ProjectNotFound ErrorType = -(1200 + iota)
	ProjectQueryError
	ProjectPermissionQueryError
	ProjectPermissionNotFound
	ProjectPermissionCannotUpdate
	ProjectPermissionExisted
	ProjectPermissionCreatingError
	ProjectCreatingError
	ProjectCannotUpdate
	ProjectCannotDelete
	ProjectRoleInvalid
)
