package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/billadm/mcp"
	"github.com/billadm/models"
)

// POST /api/v1/mcp/start
func startMcpServer(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	err := mcp.StartMcpServer()
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = mcp.GetMcpStatus()
}

// POST /api/v1/mcp/stop
func stopMcpServer(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	err := mcp.StopMcpServer()
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = mcp.GetMcpStatus()
}

// GET /api/v1/mcp/status
func getMcpStatus(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ret.Data = mcp.GetMcpStatus()
}