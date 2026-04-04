package util

import (
	"os"
	"path/filepath"
)

func GetRootDir() string {
	if Config.Mode == "debug" {
		dir, _ := os.Getwd()
		return dir
	}
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	return filepath.Dir(exePath)
}

func GetDistDir() string {
	return filepath.Join(GetRootDir(), "dist")
}

// IsDirectoryExists 判断指定路径是否为一个存在的目录
func IsDirectoryExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// IsFileExists checks if a file exists at the given path.
// Returns false if the file does not exist or if access is denied.
func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
