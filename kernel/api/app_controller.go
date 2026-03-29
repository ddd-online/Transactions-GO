package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/billadm/models"
	"github.com/billadm/workspace"
)

func exitApp(c *gin.Context) {
	ret := models.NewResult()

	// 先发送响应
	c.JSON(http.StatusOK, ret)

	// 刷新响应缓冲区，确保客户端收到响应
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	// 异步执行关闭和退出，让HTTP响应有时间发送出去
	go func() {
		logrus.Infof("--------- 退出Billadm ---------")
		workspace.Manager.Close()
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}
