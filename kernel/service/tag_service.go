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
	CreateTag(ws *workspace.Workspace, name string, categoryTransactionType string) error
	DeleteTag(ws *workspace.Workspace, ledgerId string, name string, categoryTransactionType string) error
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

func (t *tagServiceImpl) CreateTag(ws *workspace.Workspace, name string, categoryTransactionType string) error {
	logrus.Infof("start to create tag, name: %s, category: %s", name, categoryTransactionType)

	tag := &models.Tag{
		Name:                   name,
		CategoryTransactionType: categoryTransactionType,
	}

	if err := t.tagDao.CreateTag(ws, tag); err != nil {
		logrus.Errorf("create tag failed: %v", err)
		return err
	}

	logrus.Infof("create tag success, name: %s", name)
	return nil
}

func (t *tagServiceImpl) DeleteTag(ws *workspace.Workspace, ledgerId string, name string, categoryTransactionType string) error {
	logrus.Infof("start to delete tag, ledger id: %s, name: %s", ledgerId, name)

	// Delete TrTag entries that use this tag
	trTagDao := dao.GetTrTagDao()
	if err := trTagDao.DeleteTrTagByTag(ws, ledgerId, name); err != nil {
		logrus.Errorf("delete tr tags failed: %v", err)
		return err
	}

	// Delete the tag
	if err := t.tagDao.DeleteTag(ws, name, categoryTransactionType); err != nil {
		logrus.Errorf("delete tag failed: %v", err)
		return err
	}

	logrus.Infof("delete tag success, ledger id: %s, name: %s", ledgerId, name)
	return nil
}
