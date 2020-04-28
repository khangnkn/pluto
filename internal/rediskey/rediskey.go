package rediskey

import "fmt"

func AllTools() string {
	return "pluto:tool:all"
}

func ProjectByID(id uint64) string {
	return fmt.Sprintf("pluto:project:id:%d", id)
}
