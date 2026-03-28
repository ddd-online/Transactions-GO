package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/service"
	"github.com/billadm/workspace"
)

func queryCategoryByType(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	// 支持all, income, expense, transfer
	trType := c.Param("type")
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
