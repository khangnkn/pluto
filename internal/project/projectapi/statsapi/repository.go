package statsapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/logger"
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
	BuildTaskReport(projectID uint64) ([]TaskStatusPair, error)
	BuildMemberReport(projectID uint64) (MemberStatsResponse, error)
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

func (r *repository) BuildTaskReport(projectID uint64) ([]TaskStatusPair, error) {
	tasks, _, err := r.taskRepo.GetTasksByProject(projectID, task.Any, 0, 0)
	if err != nil {
		return nil, err
	}
	logger.Infof("retrieve %d tasks for project %d", len(tasks), projectID)
	var lablingCount, reviewingCount, doneCount int
	for _, t := range tasks {
		switch t.Status {
		case task.Labeling:
			lablingCount++
		case task.Reviewing:
			reviewingCount++
		case task.Done:
			doneCount++
		}
	}
	return []TaskStatusPair{
		{
			Name:  "Labling",
			Value: lablingCount,
		},
		{
			Name:  "Reviewing",
			Value: reviewingCount,
		},
		{
			Name:  "Done",
			Value: doneCount,
		},
	}, nil
}

func (r *repository) BuildMemberReport(projectID uint64) (resp MemberStatsResponse, err error) {
	labelerMap := make(map[uint64]interface{})
	reviewerMap := make(map[uint64]interface{})
	tasks, _, err := r.taskRepo.GetTasksByProject(projectID, task.Any, 0, 0)
	if err != nil {
		return
	}
	for _, v := range tasks {
		labelerMap[v.Labeler] = true
		reviewerMap[v.Reviewer] = true
	}
	resp.Labeler = len(labelerMap)
	resp.Reviewer = len(reviewerMap)
	return
}