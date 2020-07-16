package rediskey

import (
	"fmt"
)

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

func ProjectByWorkspaceID(id uint64, offset, limit int) (specKey, totalKey, pattern string) {
	specKey = fmt.Sprintf("pluto:projects:workspace:%d:offset:%d:limit%d", id, offset, limit)
	totalKey = fmt.Sprintf("pluto:projects:workspace:%d:count", id)
	pattern = fmt.Sprintf("pluto:projects:workspace:%d:*", id)
	return
}

func ProjectCountAllByWorkspaceID(workspaceID uint64) string {
	return fmt.Sprintf("pluto:projects:workspace:id:%d:count", workspaceID)
}

func ProjectByWorkspaceIDPattern(id uint64) string {
	return fmt.Sprintf("pluto:projects:workspace:id:%d:*", id)

}

func PermissionsByUserID(id uint64, role int32, offset, limit int) (string, string) {
	return fmt.Sprintf("pluto:project:permissions:userid:%d:role:%d:offset:%d:limit:%d", id, role, offset, limit),
		fmt.Sprintf("pluto:project:permissions:userid:%d:role:%d:total", id, role)
}

func ProjectPermissionByUserPattern(userID uint64) string {
	return fmt.Sprintf("pluto:project:permissions:userid:%d:*", userID)
}

func ProjectPermissionByID(projectID uint64, role uint32, offset, limit int) (specKey, totalKey, pattern string) {
	specKey = fmt.Sprintf("pluto:permissions:project:id:%d:role:%d:offset:%d:limit:%d", projectID, role, offset, limit)
	totalKey = fmt.Sprintf("pluto:permissions:project:id:%d:role:%d:total", projectID, role)
	pattern = fmt.Sprintf("pluto:permissions:project:id:%d:*", projectID)
	return
}

func TaskByID(id uint64) string {
	return fmt.Sprintf("pluto:tasks:id:%d", id)
}

func TaskByUser(userID uint64, role, status uint32, offset, limit int) (specificKey, totalKey, pattern string) {
	specificKey = fmt.Sprintf("pluto:tasks:userid:%d:role:%d:status:%d:offset:%d:limit:%d", userID, role, status, offset, limit)
	totalKey = fmt.Sprintf("pluto:tasks:userid:%d:role:%d:total", userID, role)
	pattern = fmt.Sprintf("pluto:tasks:userid:%d:*", userID)
	return
}

func TaskByProject(projectID uint64, status uint32, offset, limit int) (specificKey, totalKey, pattern string) {
	specificKey = fmt.Sprintf("pluto:tasks:project:%d:status:%d:offset:%d:limit:%d", projectID, status, offset, limit)
	totalKey = fmt.Sprintf("pluto:tasks:project:%d:status:%d:total", projectID, status)
	pattern = fmt.Sprintf("pluto:tasks:project:%d:*", projectID)
	return
}

func WorkspaceByID(id uint64) string {
	return fmt.Sprintf("pluto:workspace:id:%d", id)
}

func WorkspacesByUserID(userID uint64, role int32, offset, limit int) (string, string, string) {
	return fmt.Sprintf("pluto:workspaces:userid:%d:role:%d:offset:%d:limit:%d", userID, role, offset, limit),
		fmt.Sprintf("pluto:workspaces:userid:%d:role:%d:total", userID, role),
		fmt.Sprintf("pluto:workspaces:userid:%d:*", userID)
}

func WorkspacesByUserIDPattern(userID uint64) string {
	return fmt.Sprintf("pluto:workspaces:userid:%d:*", userID)
}

func WorkspacesPermissionByWorkspaceID(workspaceID uint64, role int32, offset, limit int) (string, string) {
	return fmt.Sprintf("pluto:workspaces:permission:workspace:%d:role:%d:offset:%d:limit:%d", workspaceID, role, offset, limit),
		fmt.Sprintf("pluto:workspaces:permission:workspace:%d:role:%d:total", workspaceID, role)
}

func WorkspacesPermissionByWorkspaceIDPattern(workspaceID uint64) string {
	return fmt.Sprintf("pluto:workspaces:permission:workspace:%d:*", workspaceID)
}
