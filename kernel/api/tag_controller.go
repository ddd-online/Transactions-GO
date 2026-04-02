package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
)

// GET /tags?categoryTransactionType=xxx
func listTags(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	categoryTransactionType := c.Query("categoryTransactionType")
	tags, err := service.GetTagService().QueryTags(ws, categoryTransactionType)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	tagDtos := make([]dto.TagDto, 0)
	for _, tag := range tags {
		tagDto := dto.TagDto{}
		tagDto.FromTag(&tag)
		tagDtos = append(tagDtos, tagDto)
	}

	ret.Data = tagDtos
}

// POST /tags
func createTag(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ret.Code = -1
		ret.Msg = "Invalid request: " + err.Error()
		return
	}

	if err := service.GetTagService().CreateTag(ws, req.Name, req.CategoryTransactionType); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

// DELETE /tags/:name
func deleteTag(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	name := c.Param("name")
	categoryTransactionType := c.Query("categoryTransactionType")
	ledgerID := c.Query("ledgerId")
	if name == "" || categoryTransactionType == "" || ledgerID == "" {
		ret.Code = -1
		ret.Msg = "Missing required parameters"
		return
	}

	if err := service.GetTagService().DeleteTag(ws, ledgerID, name, categoryTransactionType); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}