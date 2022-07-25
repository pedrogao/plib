package queue

import (
	"fmt"
	"os"
)

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fileInfo.IsDir() {
		return true
	}

	return false
}

func mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return fmt.Errorf("make dir: %s err: %s", path, err)
}
