package service

import (
	"fmt"
	"sync"

	"github.com/billadm/dao"
	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/pkg/operator"
	"github.com/billadm/util"
	"github.com/billadm/workspace"
	"github.com/sirupsen/logrus"
)

var (
	trService     TransactionRecordService
	trServiceOnce sync.Once
)

func GetTrService() TransactionRecordService {
	if trService != nil {
		return trService
	}

	trServiceOnce.Do(func() {
		trService = &transactionRecordServiceImpl{
			trDao:    dao.GetTrDao(),
			trTagDao: dao.GetTrTagDao(),
		}
	})

	return trService
}

type TransactionRecordService interface {
	CreateTr(ws *workspace.Workspace, dto *dto.TransactionRecordDto) (string, error)
	BatchCreateTr(ws *workspace.Workspace, dtos []*dto.TransactionRecordDto) (int, error)
	QueryTrsOnCondition(ws *workspace.Workspace, condition *dto.TrQueryCondition) (*dto.TrQueryResult, error)
	QueryTrsForChart(ws *workspace.Workspace, req *dto.ChartQueryRequest) (*dto.ChartQueryResponse, error)
	DeleteTrById(ws *workspace.Workspace, trId string) error
}

var _ TransactionRecordService = &transactionRecordServiceImpl{}

type transactionRecordServiceImpl struct {
	trDao    dao.TransactionRecordDao
	trTagDao dao.TrTagDao
}

// CreateTr creates a transaction record and its tags in a single transaction.
func (t *transactionRecordServiceImpl) CreateTr(ws *workspace.Workspace, trDto *dto.TransactionRecordDto) (string, error) {
	logrus.Infof("start to create transaction record, ledger id: %s, description: %s", trDto.LedgerID, trDto.Description)

	transactionID := util.GetUUID()

	record := trDto.ToTransactionRecord()
	record.TransactionID = transactionID

	// Use transaction for atomicity
	err := ws.Transaction(func(tx *workspace.Workspace) error {
		if err := t.trDao.CreateTr(tx, record); err != nil {
			return fmt.Errorf("create transaction record: %w", err)
		}

		trTags := make([]*models.TrTag, 0, len(trDto.Tags))
		for _, tag := range trDto.Tags {
			trTag := &models.TrTag{
				LedgerID:      trDto.LedgerID,
				TransactionID: transactionID,
				Tag:           tag,
			}
			trTags = append(trTags, trTag)
		}
		if err := t.trTagDao.CreateTrTags(tx, trTags); err != nil {
			return fmt.Errorf("create tr tags: %w", err)
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("create transaction record failed: %v", err)
		return "", err
	}

	logrus.Infof("create transaction record success, ledger id: %s, description: %s", trDto.LedgerID, trDto.Description)
	return transactionID, nil
}

// BatchCreateTr creates multiple transaction records in a single transaction.
func (t *transactionRecordServiceImpl) BatchCreateTr(ws *workspace.Workspace, dtos []*dto.TransactionRecordDto) (int, error) {
	logrus.Infof("start to batch create %d transaction records", len(dtos))

	if len(dtos) == 0 {
		return 0, nil
	}

	successCount := 0

	// Use transaction for atomicity
	err := ws.Transaction(func(tx *workspace.Workspace) error {
		for _, trDto := range dtos {
			transactionID := util.GetUUID()

			record := trDto.ToTransactionRecord()
			record.TransactionID = transactionID

			if err := t.trDao.CreateTr(tx, record); err != nil {
				logrus.Errorf("batch create: create transaction record failed: %v", err)
				return fmt.Errorf("create transaction record: %w", err)
			}

			trTags := make([]*models.TrTag, 0, len(trDto.Tags))
			for _, tag := range trDto.Tags {
				trTag := &models.TrTag{
					LedgerID:      trDto.LedgerID,
					TransactionID: transactionID,
					Tag:           tag,
				}
				trTags = append(trTags, trTag)
			}
			if err := t.trTagDao.CreateTrTags(tx, trTags); err != nil {
				logrus.Errorf("batch create: create tr tags failed: %v", err)
				return fmt.Errorf("create tr tags: %w", err)
			}

			successCount++
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("batch create transaction records failed: %v", err)
		return successCount, err
	}

	logrus.Infof("batch create transaction records success, count: %d", successCount)
	return successCount, nil
}

// convertSortFields converts DTO sort fields to operator sort fields
func convertSortFields(dtoSortFields []dto.QueryConditionSortField) []operator.SortField {
	if len(dtoSortFields) == 0 {
		return []operator.SortField{
			{Field: "transactionAt", Order: operator.Desc},
		}
	}
	result := make([]operator.SortField, 0, len(dtoSortFields))
	for _, sf := range dtoSortFields {
		order := operator.Desc
		if sf.Order == "asc" {
			order = operator.Asc
		}
		result = append(result, operator.SortField{
			Field: sf.Field,
			Order: order,
		})
	}
	return result
}

func (t *transactionRecordServiceImpl) QueryTrsOnCondition(ws *workspace.Workspace, condition *dto.TrQueryCondition) (*dto.TrQueryResult, error) {
	logrus.Infof("start to query trs, condition: %#v", condition)

	// Query all matching transaction records
	trs, err := t.trDao.QueryTrsOnCondition(ws, condition)
	if err != nil {
		return nil, err
	}

	// Batch query all tags in a single query (fixes N+1 problem)
	trIds := make([]string, len(trs))
	for i, tr := range trs {
		trIds[i] = tr.TransactionID
	}
	tagMap, err := t.trTagDao.QueryTrTagsByTrIds(ws, trIds)
	if err != nil {
		return nil, err
	}

	// Assemble DTOs
	trDtos := make([]*dto.TransactionRecordDto, 0, len(trs))
	for _, tr := range trs {
		trDto := &dto.TransactionRecordDto{}
		trDto.FromTransactionRecord(tr)
		if tags, ok := tagMap[tr.TransactionID]; ok {
			for _, tag := range tags {
				trDto.Tags = append(trDto.Tags, tag.Tag)
			}
		}
		trDtos = append(trDtos, trDto)
	}

	// Filter, sort, paginate and summarize
	sortFields := convertSortFields(condition.SortFields)
	summary := operator.NewTrOperator().
		Add(trDtos).
		Filter(condition.Items).
		Sort(sortFields).
		Page(condition.Offset, condition.Limit).
		Summary()

	logrus.Infof("query trs by page success, len: %d", len(summary.Items))
	return summary, nil
}

func (t *transactionRecordServiceImpl) QueryTrsForChart(ws *workspace.Workspace, req *dto.ChartQueryRequest) (*dto.ChartQueryResponse, error) {
	logrus.Infof("start to query trs for chart, granularity: %s, lines: %d", req.Granularity, len(req.Lines))

	// Collect all transaction types needed for chart lines
	ttSet := make(map[string]bool)
	for _, line := range req.Lines {
		ttSet[line.TransactionType] = true
	}

	// Query all matching transaction records from DAO (without detailed filter)
	trs, err := t.trDao.QueryTrsOnCondition(ws, &dto.TrQueryCondition{
		LedgerID: req.LedgerID,
		TsRange:  req.TsRange,
	})
	if err != nil {
		return nil, err
	}

	// Batch query tags
	trIds := make([]string, len(trs))
	for i, tr := range trs {
		trIds[i] = tr.TransactionID
	}
	tagMap, err := t.trTagDao.QueryTrTagsByTrIds(ws, trIds)
	if err != nil {
		return nil, err
	}

	// Build DTOs with tags
	trDtos := make([]*dto.TransactionRecordDto, 0, len(trs))
	for _, tr := range trs {
		trDto := &dto.TransactionRecordDto{}
		trDto.FromTransactionRecord(tr)
		if tags, ok := tagMap[tr.TransactionID]; ok {
			for _, tag := range tags {
				trDto.Tags = append(trDto.Tags, tag.Tag)
			}
		}
		trDtos = append(trDtos, trDto)
	}

	// Process each chart line - filter records per line
	response := &dto.ChartQueryResponse{
		Lines: make([]dto.ChartLineData, 0, len(req.Lines)),
	}

	for _, line := range req.Lines {
		// Filter by transaction type and outlier flag first
		var filtered []*dto.TransactionRecordDto
		for _, tr := range trDtos {
			if tr.TransactionType != line.TransactionType {
				continue
			}
			if !line.IncludeOutlier && tr.Outlier {
				continue
			}
			filtered = append(filtered, tr)
		}

		// Apply conditions with OR logic using TrOperator
		// TrOperator.Filter uses OR between conditions (any condition matches)
		filtered = operator.NewTrOperator().
			Add(filtered).
			Filter(line.Conditions).
			Summary().
			Items

		response.Lines = append(response.Lines, dto.ChartLineData{
			Label: line.Label,
			Type:  line.TransactionType,
			Items: filtered,
		})
	}

	logrus.Infof("query trs for chart success, lines: %d", len(response.Lines))
	return response, nil
}

func (t *transactionRecordServiceImpl) DeleteTrById(ws *workspace.Workspace, trId string) error {
	logrus.Infof("start to delete transaction record, tr id: %s", trId)

	// Use transaction for atomicity
	err := ws.Transaction(func(tx *workspace.Workspace) error {
		if err := t.trTagDao.DeleteTrTagByTrId(tx, trId); err != nil {
			return fmt.Errorf("delete tr tags: %w", err)
		}
		if err := t.trDao.DeleteTrById(tx, trId); err != nil {
			return fmt.Errorf("delete transaction record: %w", err)
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("delete transaction record failed: %v", err)
		return err
	}

	logrus.Infof("delete transaction record success, tr id: %s", trId)
	return nil
}
