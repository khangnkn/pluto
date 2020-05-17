package workspace

import (
	"github.com/nkhang/pluto/internal/rediskey"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
)

type Repository interface {
	Get(id uint64) (Workspace, error)
	GetByUserID(userID uint64) ([]Workspace, error)
}

type repository struct {
	dbRepo    DBRepository
	cacheRepo cache.Cache
}

func NewRepository(dbRepo DBRepository, c cache.Cache) *repository {
	return &repository{
		dbRepo:    dbRepo,
		cacheRepo: c,
	}
}

func (r *repository) Get(id uint64) (Workspace, error) {
	var w Workspace
	k := rediskey.WorkspaceByID(id)
	err := r.cacheRepo.Get(k, &w)
	if err == nil {
		return w, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for project %d", id)
	} else {
		logger.Errorf("error getting cache for workspace %d", id)
	}
	w, err = r.dbRepo.Get(id)
	if err != nil {
		logger.Error("error getting workspace from database", err)
		return Workspace{}, err
	}
	logger.Infof("getting workspace [%d] successfully", id)
	return w, nil
}

func (r *repository) GetByUserID(userID uint64) ([]Workspace, error) {
	var workspaces = make([]Workspace, 0)
	k := rediskey.WorkspacesByUserID(userID)
	err := r.cacheRepo.Get(k, &workspaces)
	if err == nil {
		return workspaces, nil
	}
	if errors.Type(err) == errors.CacheNotFound {
		logger.Infof("cache miss for user %d", userID)
	} else {
		logger.Errorf("error getting cache workspaces for user %d", userID)
	}
	workspaces, err = r.dbRepo.GetByUserID(userID)
	if err != nil {
		logger.Error("error getting workspaces from database", err)
		return nil, err
	}
	logger.Infof("getting workspace for user %d successfully", userID)
	return workspaces, nil
}
