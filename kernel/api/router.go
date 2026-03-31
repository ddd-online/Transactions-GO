package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/billadm/models"
)

func ServeAPI(ginServer *gin.Engine) {
	v1 := ginServer.Group("/api/v1")
	{
		// App control
		v1.POST("/app/exit", exitApp)

		// Ledgers: RESTful CRUD
		ledgers := v1.Group("/ledgers")
		{
			ledgers.GET("", listLedgers)
			ledgers.POST("", createLedger)
			ledgers.GET("/:id", getLedger)
			ledgers.PATCH("/:id", updateLedger)
			ledgers.DELETE("/:id", deleteLedger)
		}

		// Transactions: query uses POST for complex filters, others RESTful
		transactions := v1.Group("/transactions")
		{
			transactions.POST("/query", queryTransactions)
			transactions.POST("/batch", batchCreateTransactions)
			transactions.POST("", createTransaction)
			transactions.DELETE("/:id", deleteTransaction)
		}

		// Categories: GET by type query param
		v1.GET("/categories", listCategories)

		// Tags: GET by category query param
		v1.GET("/tags", listTags)

		// Workspace
		workspace := v1.Group("/workspace")
		{
			workspace.POST("", openWorkspace)
			workspace.GET("/status", getWorkspaceStatus)
		}
	}
}

func JsonArg(c *gin.Context, result *models.Result) (arg map[string]any, ok bool) {
	arg = make(map[string]any)
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = fmt.Sprintf("parses request failed: %v", err)
		return
	}

	ok = true
	return
}

func PathParam(c *gin.Context, name string) (val string, ok bool) {
	val = c.Param(name)
	ok = val != ""
	return
}
