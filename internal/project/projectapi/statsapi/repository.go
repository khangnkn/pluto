package statsapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/task"
)

type repository struct {
	datasetRepo dataset.Repository
	taskRepo    task.Repository
	imageRepo   image.Repository
}

func NewRepository(d dataset.Repository, t task.Repository, i image.Repository) *repository {
	return &repository{
		datasetRepo: d,
		taskRepo:    t,
		imageRepo:   i,
	}
}

type Repository interface {
	BuildReport(datasetID uint64) (DatasetStatsResponse, error)
}

func (r *repository) BuildReport(datasetID uint64) (DatasetStatsResponse, error) {
	details, err := r.imageRepo.GetAllImageByDataset(datasetID)
	if err != nil {
		return DatasetStatsResponse{}, err
	}
	return DatasetStatsResponse{
		AnnotatedTimes:      buildAnnotatedTimesPair(details),
		AnnotatedStatusPair: buildAnnotatedStatusPairs(details),
	}, nil
}

func buildAnnotatedTimesPair(images []image.Image) (result []AnnotatedTimePair) {
	var resMap = make(map[uint32]int)
	for _, i := range images {
		times := i.Status
		resMap[times] = resMap[times] + 1
	}
	for k, v := range resMap {
		result = append(result, AnnotatedTimePair{
			Times: k,
			Count: v,
		})
	}
	return
}

func buildAnnotatedStatusPairs(images []image.Image) []AnnotatedStatusPair {
	var notAnnotatedCount, annotatedCount int
	for _, v := range images {
		if v.Status == 0 {
			notAnnotatedCount++
		} else {
			annotatedCount++
		}

	}
	return []AnnotatedStatusPair{
		{
			Name:  "Unused",
			Value: notAnnotatedCount,
		},
		{
			Name:  "Finished",
			Value: annotatedCount,
		},
	}
}
