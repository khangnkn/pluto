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

func ProjectByID(wID, pID uint64) string {
	return fmt.Sprintf("pluto:project:wid:%d:pid:%d", wID, pID)
}

func ProjectByWorkspaceID(id uint64) string {
	return fmt.Sprintf("pluto:projects:workspace:id:%d", id)
}

func WorkspaceByID(id uint64) string {
	return fmt.Sprintf("pluto:workspace:id:%d", id)
}
