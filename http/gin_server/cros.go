package gin_server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewCorsMW(config *CORSConfig) gin.HandlerFunc {
	crosConfig := cors.DefaultConfig()
	if len(config.AllowOrigins) != 0 {
		crosConfig.AllowOrigins = config.AllowOrigins
	}
	crosConfig.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	crosConfig.ExposeHeaders = []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"}
	return cors.New(crosConfig)
}
