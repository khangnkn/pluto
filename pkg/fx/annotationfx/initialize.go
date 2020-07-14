package annotationfx

import (
	"github.com/nats-io/nats.go"
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/annotation"
)

func provideAnnotationService(workspaceRepo workspace.Repository,
	projectRepo project.Repository,
	datasetRepo dataset.Repository,
	labelRepo label.Repository,
	nc *nats.EncodedConn) annotation.Service {
	return annotation.NewService(workspaceRepo, projectRepo, datasetRepo, labelRepo, nc)
}
