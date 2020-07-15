package errors

const (
	ProjectNotFound ErrorType = -(1200 + iota)
	ProjectQueryError
	ProjectPermissionQueryError
	ProjectPermissionNotFound
	ProjectPermissionCannotUpdate
	ProjectPermissionExisted
	ProjectPermissionCreatingError
	ProjectPermissionCannotDelete
	ProjectCreatingError
	ProjectCannotUpdate
	ProjectCannotDelete
	ProjectRoleInvalid
)
