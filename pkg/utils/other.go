package utils

import (
	"fmt"
)

func SizeToHuman(space int64) string {
	if space < 1024 {
		return fmt.Sprintf("%d kB", space)
	}
	if space < 1024*1024 {
		return fmt.Sprintf("%.01f MB", float32(space)/1024)
	}
	if space < 1024*1024*1024 {
		return fmt.Sprintf("%.01f GB", float32(space)/1024/1024)
	}
	return fmt.Sprintf("%.02f TB", float32(space)/1024/1024/1024)
}
