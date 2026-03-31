package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
)

// POST /transactions/query
func queryTransactions(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	queryCondition, ok := dto.JsonQueryCondition(c, ret)
	if !ok {
		return
	}
	logrus.Debugf("query condition: %v", queryCondition)

	result, err := service.GetTrService().QueryTrsOnCondition(ws, queryCondition)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = result
}

// POST /transactions
func createTransaction(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	trDto, ok := dto.JsonTransactionRecordDto(c, ret)
	if !ok {
		return
	}
	logrus.Debugf("tr dto: %v", trDto)

	if !trDto.Validate(ret) {
		logrus.Errorf("validate transaction record error, err: %v", ret.Msg)
		return
	}

	trId, err := service.GetTrService().CreateTr(ws, trDto)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = trId
}

// POST /transactions/batch
func batchCreateTransactions(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	dtos, ok := dto.JsonTransactionRecordDtoBatch(c, ret)
	if !ok {
		return
	}
	logrus.Debugf("batch create %d transaction records", len(dtos))

	// Validate all records first
	for i, trDto := range dtos {
		if !trDto.Validate(ret) {
			ret.Code = -1
			ret.Msg = fmt.Sprintf("record %d: %s", i+1, ret.Msg)
			logrus.Errorf("validate transaction record %d error: %v", i+1, ret.Msg)
			return
		}
	}

	count, err := service.GetTrService().BatchCreateTr(ws, dtos)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = count
}

// DELETE /transactions/:id
func deleteTransaction(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	id := c.Param("id")
	if id == "" {
		ret.Code = -1
		ret.Msg = "missing transaction id"
		return
	}

	err := service.GetTrService().DeleteTrById(ws, id)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}
