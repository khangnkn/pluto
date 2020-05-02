package rediskey

import "fmt"

func AllTools() string {
	return "pluto:tool:all"
}

func ProjectByID(wID, pID uint64) string {
	return fmt.Sprintf("pluto:project:wid:%d:pid:%d", wID, pID)
}

func ProjectByWorkspaceID(id uint64) string {
	return fmt.Sprintf("pluto:projects:workspace:id:%d", id)
}

func LabelsByProject(pID uint64) string {
	return fmt.Sprintf("pluto:labels:project:id:%d", pID)
}

func WorkspaceByID(id uint64) string {
	return fmt.Sprintf("pluto:workspace:id:%d", id)
}
