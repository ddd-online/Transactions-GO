package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

func openWorkspace(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := JsonArg(c, ret)
	if !ok {
		return
	}

	workspaceDir, ok := arg["workspaceDir"].(string)
	if !ok || workspaceDir == "" {
		ret.Code = -1
		ret.Msg = "工作目录路径不能为空"
		return
	}

	err := workspace.Manager.OpenWorkspace(workspaceDir)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

func hasOpenedWorkspace(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	if workspace.Manager.OpenedWorkspace() == nil {
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		ret.Data = ""
		return
	}

	ret.Data = workspace.Manager.OpenedWorkspace().GetDirectory()
}
