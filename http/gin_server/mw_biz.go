package gin_server

import (
	"github.com/gin-gonic/gin"
	"github.com/ragpanda/go-toolkit/biz"
	"github.com/ragpanda/go-toolkit/log/logkit"
)

const LogIDKey = "X-Log-Id"

func BizDataMw(c *gin.Context) {
	d := c.Value(biz.BizDataKey)
	if d == nil {
		d = biz.NewBizData()
		biz.SetBizDataToGinCtx(c, d.(*biz.BizData))
	}

	bizData := d.(*biz.BizData)

	bizData.UserID = c.GetString("user")
	bizData.FromIP = c.ClientIP()

	logID := c.Request.Header.Get(LogIDKey)
	if logID == "" {
		logID = logkit.GetLogID()
		c.Request.Header.Set(LogIDKey, logID)
	}
	bizData.LogID = logID

	c.Next()
}
