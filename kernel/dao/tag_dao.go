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
