package rediskey

import "fmt"

func AllTools() string {
	return "pluto:tool:all"
}

func DatasetByID(dID uint64) string {
	return fmt.Sprintf("pluto:dataset:id:%d", dID)
}

func DatasetByProject(pID uint64) string {
	return fmt.Sprintf("pluto:dataset:project:id:%d", pID)
}

func LabelsByProject(pID uint64) string {
	return fmt.Sprintf("pluto:labels:project:id:%d", pID)
}

func ImageByID(id uint64) string {
	return fmt.Sprintf("pluto:image:id:%d", id)
}

func ImageByDatasetID(dID uint64, offset, limit int) string {
	return fmt.Sprintf("pluto:image:dataset:id:%d:offset:%d:limit:%d", dID, offset, limit)
}

func ImageAllByDatasetID(dID uint64) string {
	return fmt.Sprintf("pluto:image:dataset:id:%d", dID)
}

func ImageByDatasetIDAllKeys(dID uint64) string {
	return fmt.Sprintf("pluto:image:dataset:id:%d:*", dID)
}

func ProjectByID(pID uint64) string {
	return fmt.Sprintf("pluto:project:id:%d", pID)
}

func ProjectByWorkspaceID(id uint64) string {
	return fmt.Sprintf("pluto:projects:workspace:id:%d", id)
}

func WorkspaceByID(id uint64) string {
	return fmt.Sprintf("pluto:workspace:id:%d", id)
}

func WorkspacesByUserID(userID uint64) string {
	return fmt.Sprintf("pluto:workspaces:user:id:%d", userID)
}
