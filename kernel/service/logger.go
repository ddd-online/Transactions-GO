package service

import (
	"github.com/sirupsen/logrus"
)

// Logger abstracts logging operations for the service layer.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// workspaceLoggerAdapter adapts workspace's *logrus.Logger to the Logger interface.
type workspaceLoggerAdapter struct {
	logger *logrus.Logger
}

func (l *workspaceLoggerAdapter) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *workspaceLoggerAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *workspaceLoggerAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// WorkspaceLogger returns a Logger adapter for the workspace's logrus.Logger.
func WorkspaceLogger(ws interface{ GetLogger() *logrus.Logger }) Logger {
	return &workspaceLoggerAdapter{logger: ws.GetLogger()}
}
