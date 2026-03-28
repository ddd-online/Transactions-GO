package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
)

func queryTrOnCondition(c *gin.Context) {
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

func createTransactionRecord(c *gin.Context) {
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

func deleteTransactionRecordById(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	arg, ok := JsonArg(c, ret)
	if !ok {
		return
	}

	trId, ok := arg["trId"].(string)
	if !ok {
		ret.Code = -1
		ret.Msg = "trId在请求体中不存在"
		return
	}

	err := service.GetTrService().DeleteTrById(ws, trId)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}
