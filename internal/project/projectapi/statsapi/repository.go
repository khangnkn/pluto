package statsapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/pkg/annotation"
	"github.com/nkhang/pluto/pkg/logger"
)

type repository struct {
	datasetRepo       dataset.Repository
	taskRepo          task.Repository
	imageRepo         image.Repository
	labelRepo         label.Repository
	annotationService annotation.Service
}

func NewRepository(d dataset.Repository, t task.Repository, i image.Repository, s annotation.Service, l label.Repository) *repository {
	return &repository{
		datasetRepo:       d,
		taskRepo:          t,
		imageRepo:         i,
		annotationService: s,
		labelRepo:         l,
	}
}

type Repository interface {
	BuildReport(projectID, datasetID uint64) (DatasetStatsResponse, error)
	BuildTaskReport(projectID uint64) ([]TaskStatusPair, error)
	BuildMemberReport(projectID uint64) (MemberStatsResponse, error)
	BuildLabelReport(projectID, labelID uint64) (GetLabelStatsResponse, error)
}

func (r *repository) BuildReport(projectId, datasetID uint64) (DatasetStatsResponse, error) {
	var (
		images []image.Image
		err    error
	)
	if datasetID == 0 {
		images, err = r.getAllImagesForProject(projectId)
	} else {
		images, err = r.imageRepo.GetAllImageByDataset(datasetID)
	}
	if err != nil {
		return DatasetStatsResponse{}, err
	}
	return DatasetStatsResponse{
		AnnotatedTimes:      buildAnnotatedTimesPair(images),
		AnnotatedStatusPair: buildAnnotatedStatusPairs(images),
	}, nil
}

func (r *repository) getAllImagesForProject(projectId uint64) (resp []image.Image, err error) {
	datasets, err := r.datasetRepo.GetByProject(projectId)
	if err != nil {
		return nil, err
	}
	for i := range datasets {
		details, err := r.imageRepo.GetAllImageByDataset(datasets[i].ID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, details...)
	}
	return
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
	var labelingCount, reviewingCount, doneCount int
	for _, t := range tasks {
		switch t.Status {
		case task.Labeling:
			labelingCount++
		case task.Reviewing:
			reviewingCount++
		case task.Done:
			doneCount++
		}
	}
	return []TaskStatusPair{
		{
			Name:  "Labling",
			Value: labelingCount,
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
	resp.Labelers = len(labelerMap)
	resp.Reviewers = len(reviewerMap)
	return
}

func (r *repository) BuildLabelReport(projectID, labelID uint64) (GetLabelStatsResponse, error) {
	datasets, err := r.datasetRepo.GetByProject(projectID)
	if err != nil {
		return GetLabelStatsResponse{}, err
	}
	var images []image.Image
	for i := range datasets {
		imgs, err := r.imageRepo.GetAllImageByDataset(datasets[i].ID)
		if err != nil {
			return GetLabelStatsResponse{}, err
		}
		images = append(images, imgs...)
	}
	if labelID == 0 {
		return r.buildAllLabelReport(projectID, images)
	}
	stats, err := r.annotationService.GetLabelCount(projectID, labelID)
	if err != nil {
		return GetLabelStatsResponse{}, err
	}
	return GetLabelStatsResponse{
		TotalObjects: stats.TotalObject,
		Donut: []DonutPart{
			{
				Name:  "Have the label",
				Value: stats.TotalImage,
			},
			{
				Name:  "Don't have the label",
				Value: len(images) - stats.TotalImage,
			},
		},
	}, nil
}

func (r *repository) buildAllLabelReport(projectID uint64, images []image.Image) (resp GetLabelStatsResponse, err error) {
	labels, err := r.labelRepo.GetByProjectId(projectID)
	if err != nil {
		return
	}
	var totalObject int
	for i := range labels {
		stats, err := r.annotationService.GetLabelCount(projectID, labels[i].ID)
		if err != nil {
			logger.Errorf("error getting statistic from annotation service. project %d, label %d err %v", projectID, labels[i].ID, err)
			continue
		}
		totalObject += stats.TotalObject
	}
	var labeledImages = make([]image.Image, 0)
	for _, i := range images {
		if i.Status != 0 {
			labeledImages = append(labeledImages, i)
		}
	}
	return GetLabelStatsResponse{
		TotalObjects: totalObject,
		Donut: []DonutPart{
			{
				Name:  "Have the label",
				Value: len(labeledImages),
			},
			{
				Name:  "Don't have the label",
				Value: len(images) - len(labeledImages),
			},
		},
	}, nil
}
