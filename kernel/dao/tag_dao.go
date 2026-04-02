package dao

import (
	"sync"

	"github.com/billadm/constant"
	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

var (
	tagDao     TagDao
	tagDaoOnce sync.Once
)

func GetTagDao() TagDao {
	if tagDao != nil {
		return tagDao
	}
	tagDaoOnce.Do(func() {
		tagDao = &TagDaoImpl{}
	})
	return tagDao
}

type TagDao interface {
	QueryTags(ws *workspace.Workspace, categoryTransactionType string) ([]models.Tag, error)
	CreateTag(ws *workspace.Workspace, tag *models.Tag) error
	DeleteTag(ws *workspace.Workspace, name string, categoryTransactionType string) error
	DeleteTagsByCategory(ws *workspace.Workspace, categoryTransactionType string) error
}

var _ TagDao = &TagDaoImpl{}

type TagDaoImpl struct{}

func (t *TagDaoImpl) QueryTags(ws *workspace.Workspace, categoryTransactionType string) ([]models.Tag, error) {
	tags := make([]models.Tag, 0)
	db := ws.GetDb()
	if categoryTransactionType != constant.All {
		db = db.Where("category_transaction_type = ?", categoryTransactionType)
	}
	if err := db.Order("name DESC").Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *TagDaoImpl) CreateTag(ws *workspace.Workspace, tag *models.Tag) error {
	if err := ws.GetDb().Create(tag).Error; err != nil {
		return err
	}
	return nil
}

func (t *TagDaoImpl) DeleteTag(ws *workspace.Workspace, name string, categoryTransactionType string) error {
	if err := ws.GetDb().
		Where("name = ? AND category_transaction_type = ?", name, categoryTransactionType).
		Delete(&models.Tag{}).Error; err != nil {
		return err
	}
	return nil
}

func (t *TagDaoImpl) DeleteTagsByCategory(ws *workspace.Workspace, categoryTransactionType string) error {
	if err := ws.GetDb().
		Where("category_transaction_type = ?", categoryTransactionType).
		Delete(&models.Tag{}).Error; err != nil {
		return err
	}
	return nil
}
