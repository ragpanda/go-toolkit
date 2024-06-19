package gin_server

import (
	"github.com/gin-gonic/gin"
	"github.com/ragpanda/go-toolkit/log/logkit"
)

const LogIDKey = "logid"

func LogIDMW(c *gin.Context) {
	logID := c.Request.Header.Get(LogIDKey)
	if logID == "" {
		logID = logkit.GetLogID()
		c.Request.Header.Set(LogIDKey, logID)
	}

	GetBizData(c).SetKey(LogIDKey, logID)
	c.Next()
}
