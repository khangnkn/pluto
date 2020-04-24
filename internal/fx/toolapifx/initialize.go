package toolapifx

import "github.com/nkhang/pluto/internal/toolapi"

func initializer(r toolapi.Repository) toolapi.ToolRepository {
	return toolapi.NewRepository(r)
}
