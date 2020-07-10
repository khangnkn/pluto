package workspacefx

import "go.uber.org/fx"

var Module = fx.Provide(
	provideWorkspaceDBRepository,
	provideWorkspaceRepository,
	provideWorkspaceAPIRepository,
	provideWorkspacePermAPIRepository,
	fx.Annotated{
		Name:   "WorkspaceService",
		Target: provideWorkspaceService,
	})
