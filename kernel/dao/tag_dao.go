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
	UpdateTagSort(ws *workspace.Workspace, name string, categoryTransactionType string, sortOrder int) error
	GetMaxSortOrder(ws *workspace.Workspace, categoryTransactionType string) (int, error)
	CountRecordsByTag(ws *workspace.Workspace, ledgerId string, tag string) (int64, error)
}

var _ TagDao = &TagDaoImpl{}

type TagDaoImpl struct{}

func (t *TagDaoImpl) QueryTags(ws *workspace.Workspace, categoryTransactionType string) ([]models.Tag, error) {
	tags := make([]models.Tag, 0)
	db := ws.GetDb()
	if categoryTransactionType != constant.All {
		db = db.Where("category_transaction_type = ?", categoryTransactionType)
	}
	if err := db.Order("sort_order ASC, name DESC").Find(&tags).Error; err != nil {
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

func (t *TagDaoImpl) UpdateTagSort(ws *workspace.Workspace, name string, categoryTransactionType string, sortOrder int) error {
	if err := ws.GetDb().
		Model(&models.Tag{}).
		Where("name = ? AND category_transaction_type = ?", name, categoryTransactionType).
		Update("sort_order", sortOrder).Error; err != nil {
		return err
	}
	return nil
}

func (t *TagDaoImpl) GetMaxSortOrder(ws *workspace.Workspace, categoryTransactionType string) (int, error) {
	var maxSortOrder int
	err := ws.GetDb().Model(&models.Tag{}).
		Where("category_transaction_type = ?", categoryTransactionType).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error
	if err != nil {
		return 0, err
	}
	return maxSortOrder, nil
}

func (t *TagDaoImpl) CountRecordsByTag(ws *workspace.Workspace, ledgerId string, tag string) (int64, error) {
	var count int64
	err := ws.GetDb().Model(&models.TrTag{}).
		Where("ledger_id = ? AND tag = ?", ledgerId, tag).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}