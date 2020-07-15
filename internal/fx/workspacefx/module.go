package workspacefx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideWorkspaceDBRepository,
	provideWorkspaceRepository,
	provideWorkspaceAPIRepository,
	fx.Annotated{
		Name:   "WorkspaceService",
		Target: provideWorkspaceService,
	})
