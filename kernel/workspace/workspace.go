package workspace

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/billadm/constant"
	"github.com/billadm/logger"
	"github.com/billadm/util"
)

type Workspace struct {
	directory string
	logger    *logrus.Logger
	db        *gorm.DB
}

func NewWorkspace(directory string) (*Workspace, error) {
	if !util.IsDirectoryExists(directory) {
		err := os.MkdirAll(directory, 0750)
		if err != nil {
			return nil, err
		}
	}
	// 实例化log
	log := logrus.New()
	logFile := filepath.Join(directory, constant.LogName)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	log.SetFormatter(&logger.CustomFormatter{})
	logLevel, err := logrus.ParseLevel(util.Config.LogLevel)
	if err != nil {
		return nil, err
	}
	log.SetLevel(logLevel)
	// 实例化db
	dbFile := filepath.Join(directory, constant.DbName)
	if err := util.OpenAndInit(dbFile); err != nil {
		return nil, err
	}
	db, err := util.NewDbInstance(dbFile)
	if err != nil {
		return nil, err
	}

	return &Workspace{
		directory: directory,
		logger:    log,
		db:        db,
	}, nil
}

func (w *Workspace) GetLogger() *logrus.Logger {
	return w.logger
}

func (w *Workspace) GetDb() *gorm.DB {
	return w.db
}

func (w *Workspace) GetDirectory() string {
	return w.directory
}

func (w *Workspace) Close() {
	sqlDb, err := w.db.DB()
	if err != nil {
		w.logger.Error(err)
	}
	err = sqlDb.Close()
	if err != nil {
		w.logger.Error(err)
	}
}
