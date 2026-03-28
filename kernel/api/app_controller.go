package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

func exitApp(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	logrus.Infof("--------- 退出Billadm ---------")
	workspace.Manager.Close()
	os.Exit(0)
}
