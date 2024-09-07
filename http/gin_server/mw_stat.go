package gin_server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ragpanda/go-toolkit/log"
	"github.com/ragpanda/go-toolkit/metrics"
)

func StatMW(c *gin.Context) {
	start := time.Now()
	defer func() {
		labels := metrics.APIRequest{
			Method:           c.Request.Method,
			Path:             c.Request.URL.Path,
			StatusCode:       c.Request.Response.StatusCode,
			Success:          false,
			Duration:         0,
			BusinessCategory: "",
		}
		end := time.Now()
		labels.Duration = end.Sub(start)
		if 200 < c.Request.Response.StatusCode && c.Request.Response.StatusCode < 400 {
			labels.Success = true
		}
		metrics.RecordAPIRequest(labels)
		log.Info(c, "[stat] api:%s %s %d %v reqs:%d,resps:%d", labels.Method, labels.Path, labels.StatusCode, labels.Duration,
			c.Request.ContentLength, c.Request.Response.ContentLength)
	}()
	c.Next()
}
