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

// POST /templates
func createTemplate(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	templateDto, ok := dto.JsonTransactionTemplateDto(c, ret)
	if !ok {
		return
	}
	logrus.Debugf("template dto: %v", templateDto)

	if !templateDto.Validate(ret) {
		logrus.Errorf("validate transaction template error, err: %v", ret.Msg)
		return
	}

	templateId, err := service.GetTrTemplateService().Create(ws, templateDto)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = templateId
}

// GET /templates
func listTemplates(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	ledgerId := c.Query("ledgerId")
	if ledgerId == "" {
		ret.Code = -1
		ret.Msg = "missing ledgerId"
		return
	}

	templates, err := service.GetTrTemplateService().ListByLedgerId(ws, ledgerId)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = templates
}

// DELETE /templates/:id
func deleteTemplate(c *gin.Context) {
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
		ret.Msg = "missing template id"
		return
	}

	err := service.GetTrTemplateService().DeleteById(ws, id)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}