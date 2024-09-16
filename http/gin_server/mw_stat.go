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
			Success:          false,
			Duration:         0,
			BusinessCategory: "",
		}
		end := time.Now()
		labels.Duration = end.Sub(start)

		if c.Request.Response != nil {
			labels.StatusCode = c.Request.Response.StatusCode
			if 200 < c.Request.Response.StatusCode && c.Request.Response.StatusCode < 400 {
				labels.Success = true
			}
		} else {
			labels.StatusCode = 1
			labels.Success = false
		}

		metrics.RecordAPIRequest(labels)
		log.Info(c, "[stat] api:%s %s %d %v reqs:%d, log_id:%s",
			labels.Method, labels.Path,
			labels.StatusCode, labels.Duration,
			c.Request.ContentLength,
			c.Request.Header.Get("X-Log-Id"),
		)
	}()
	c.Next()
}
