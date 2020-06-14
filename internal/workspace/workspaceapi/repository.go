package workspaceapi

import "github.com/nkhang/pluto/internal/workspace"

type Repository interface {
	GetByID(id uint64) (WorkspaceInfoResponse, error)
	GetByUserID(userID uint64) ([]WorkspaceInfoResponse, error)
	CreateWorkspace(p CreateWorkspaceRequest) error
}

type repository struct {
	workspaceRepository workspace.Repository
}

func NewRepository(r workspace.Repository) *repository {
	return &repository{r}
}

func (r *repository) GetByID(id uint64) (WorkspaceInfoResponse, error) {
	w, err := r.workspaceRepository.Get(id)
	if err != nil {
		return WorkspaceInfoResponse{}, err
	}
	return toWorkspaceInfoResponse(w), nil
}

func (r *repository) GetByUserID(userID uint64) ([]WorkspaceInfoResponse, error) {
	w, err := r.workspaceRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]WorkspaceInfoResponse, len(w))
	for i := range w {
		responses[i] = toWorkspaceInfoResponse(w[i])
	}
	return responses, nil
}

func (r *repository) CreateWorkspace(p CreateWorkspaceRequest) error {
	return r.workspaceRepository.Create(p.UserID, p.Title, p.Description)
}
