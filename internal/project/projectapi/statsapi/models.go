package statsapi

type GetDatasetStatsRequest struct {
	DatasetID uint64 `form:"dataset_id" json:"dataset_id"`
}

type DatasetStatsResponse struct {
	AnnotatedTimes      []AnnotatedTimePair   `json:"annotated_times"`
	AnnotatedStatusPair []AnnotatedStatusPair `json:"annotated_status"`
}

type MemberStatsResponse struct {
	Labeler  int `json:"labeler"`
	Reviewer int `json:"reviewer"`
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
