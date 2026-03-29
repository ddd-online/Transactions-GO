package service

import (
	"sync"

	"github.com/billadm/dao"
	"github.com/billadm/models"
	"github.com/billadm/workspace"
	"github.com/sirupsen/logrus"
)

var (
	tagService     TagService
	tagServiceOnce sync.Once
)

func GetTagService() TagService {
	if tagService != nil {
		return tagService
	}

	tagServiceOnce.Do(func() {
		tagService = &tagServiceImpl{
			tagDao: dao.GetTagDao(),
		}
	})

	return tagService
}

type TagService interface {
	QueryTags(ws *workspace.Workspace, categoryTransactionType string) ([]models.Tag, error)
}

var _ TagService = &tagServiceImpl{}

type tagServiceImpl struct {
	tagDao dao.TagDao
}

func (t *tagServiceImpl) QueryTags(ws *workspace.Workspace, categoryTransactionType string) ([]models.Tag, error) {
	logrus.Info("start to query tag")
	tags, err := t.tagDao.QueryTags(ws, categoryTransactionType)
	if err != nil {
		return nil, err
	}

	logrus.Infof("query tag success, length: %v", tags)
	return tags, nil
}
