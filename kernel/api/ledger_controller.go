package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/billadm/constant"
	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
)

func queryAllLedgers(c *gin.Context) {
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

	ledgerId, ok := arg["id"].(string)
	if !ok {
		ret.Code = -1
		ret.Msg = "'id'在请求体中不存在"
		return
	}

	var ledgers []models.Ledger
	var err error
	if ledgerId == constant.All {
		ledgers, err = service.GetLedgerService().ListAllLedger(ws)
		if err != nil {
			ret.Code = -1
			ret.Msg = err.Error()
			return
		}
	} else {
		var ledger *models.Ledger
		ledgerIds := strings.Split(ledgerId, ",")
		for _, id := range ledgerIds {
			id = strings.TrimSpace(id)
			ledger, err = service.GetLedgerService().QueryLedgerById(ws, id)
			if err != nil {
				ret.Code = -1
				ret.Msg = fmt.Sprintf("查询账本 %s 失败: %v", id, err)
				return
			}
			ledgers = append(ledgers, *ledger)
		}
	}

	ledgerDtos := make([]dto.LedgerDto, 0)
	for _, ledger := range ledgers {
		ledgerDto := dto.LedgerDto{}
		ledgerDto.FromLedger(&ledger)
		ledgerDtos = append(ledgerDtos, ledgerDto)
	}

	ret.Data = ledgerDtos
}

func createLedger(c *gin.Context) {
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

	ledgerName, ok := arg["name"].(string)
	if !ok {
		ret.Code = -1
		ret.Msg = "name在请求体中不存在"
		return
	}

	ledgerId, err := service.GetLedgerService().CreateLedger(ws, ledgerName)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = ledgerId
}

func modifyLedgerName(c *gin.Context) {
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

	ledgerId, ok := arg["id"].(string)
	if !ok {
		ret.Code = -1
		ret.Msg = "id在请求体中不存在"
		return
	}

	ledgerName, ok := arg["name"].(string)
	if !ok {
		ret.Code = -1
		ret.Msg = "name在请求体中不存在"
		return
	}

	err := service.GetLedgerService().ModifyLedgerName(ws, ledgerId, ledgerName)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

func deleteLedger(c *gin.Context) {
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

	trId, ok := arg["id"].(string)
	if !ok {
		ret.Code = -1
		ret.Msg = "id在请求体中不存在"
		return
	}

	err := service.GetLedgerService().DeleteLedgerById(ws, trId)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}
