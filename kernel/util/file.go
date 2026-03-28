package util

import (
	"io"
	"os"
)

func WriteStringToFile(path, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.WriteString(file, content)
	return err
}
