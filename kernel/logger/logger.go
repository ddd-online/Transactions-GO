// Package logger
package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.StandardLogger().SetOutput(os.Stdout)
	logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	logrus.StandardLogger().SetFormatter(&CustomFormatter{})
}

// Init 初始化日志配置
func Init(level string) error {
	// 设置日志级别
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.StandardLogger().SetLevel(logLevel)
	return nil
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var logLine string
	if entry.HasCaller() {
		logLine = fmt.Sprintf("[%s] [%s] [goroutine-%d] [%s:%d] %s\n", timestamp, entry.Level.String(), getGoID(), entry.Caller.File, entry.Caller.Line, entry.Message)
	} else {
		logLine = fmt.Sprintf("[%s] [%s] [goroutine-%d] %s\n", timestamp, entry.Level.String(), getGoID(), entry.Message)
	}
	return []byte(logLine), nil
}

func getGoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
