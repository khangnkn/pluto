package statsapi

type GetDatasetStatsRequest struct {
	DatasetID uint64 `form:"dataset_id" json:"dataset_id"`
}

type GetLabelStatsRequest struct {
	LabelID uint64 `form:"label_id" json:"label_id"`
}

type GetLabelStatsResponse struct {
	TotalObjects int         `json:"total_objects"`
	Donut        []DonutPart `json:"donut"`
}

type DonutPart struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type DatasetStatsResponse struct {
	AnnotatedTimes      []AnnotatedTimePair   `json:"annotated_times"`
	AnnotatedStatusPair []AnnotatedStatusPair `json:"annotated_status"`
}

type MemberStatsResponse struct {
	Labelers  int `json:"labellers"`
	Reviewers int `json:"reviewers"`
}

type AnnotatedTimePair struct {
	Times uint32 `json:"times"`
	Count int    `json:"count"`
}

type AnnotatedStatusPair struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type TaskStatusPair struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
