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
	// 获取路径的文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		// 如果是路径不存在的错误，返回 false, nil
		if os.IsNotExist(err) {
			return false
		}
		// 如果是其他错误（如权限问题），返回 false 和错误信息
		return false
	}

	// 检查是否为目录
	return fileInfo.IsDir()
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		// 如果错误是 "不存在"，则返回 false
		if os.IsNotExist(err) {
			return false
		}
	}
	// 其他情况（如权限错误等）也视为“不存在”或不可访问
	// 如果你希望区分错误类型，可以返回 (bool, error)，见下方增强版
	return true
}
