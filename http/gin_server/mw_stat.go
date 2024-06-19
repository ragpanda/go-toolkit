package gin_server

import (
	"time"

	"github.com/gin-gonic/gin"
)

func StatMW(c *gin.Context) {
	start := time.Now()
	var err error
	defer func() {
		cost := time.Now().Sub(start)
		consts.MetricAPICost.EmitTimer(cost.Milliseconds(),
			consts.MetricAPICostTag_Path.Value(ctx.Request.URL.Path))

		consts.MetricAPIQPS.EmitRate(1,
			consts.MetricAPIQPSTag_Path.Value(ctx.Request.URL.Path))
		if err != nil {
			consts.MetricAPIError.EmitCounter(1,
				consts.MetricAPIErrorTag_Path.Value(ctx.Request.URL.Path))
		}
	}()
	c.Next()
}
