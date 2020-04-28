package label

type repository struct {
	disk *diskRepository
}

func NewRepository(d *diskRepository) *repository {
	return &repository{disk: d}
}

func (d *repository) GetByProjectId(projectId uint64) ([]Label, error) {
	return d.disk.GetByProjectID(projectId)
}
