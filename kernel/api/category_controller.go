package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
)

// GET /categories?type=all|income|expense|transfer
func listCategories(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	trType := c.Query("type")
	categories, err := service.GetCategoryService().QueryCategory(ws, trType)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	categoryDtos := make([]dto.CategoryDto, 0)
	for _, category := range categories {
		categoryDto := dto.CategoryDto{}
		categoryDto.FromCategory(&category)
		categoryDtos = append(categoryDtos, categoryDto)
	}

	ret.Data = categoryDtos
}

// POST /categories
func createCategory(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ret.Code = -1
		ret.Msg = "Invalid request: " + err.Error()
		return
	}

	if err := service.GetCategoryService().CreateCategory(ws, req.LedgerID, req.Name, req.TransactionType); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

// DELETE /categories/:name
func deleteCategory(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	name := c.Param("name")
	transactionType := c.Query("type")
	ledgerID := c.Query("ledgerId")
	if name == "" || transactionType == "" || ledgerID == "" {
		ret.Code = -1
		ret.Msg = "Missing required parameters"
		return
	}

	if err := service.GetCategoryService().DeleteCategory(ws, ledgerID, name, transactionType); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}