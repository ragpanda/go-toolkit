package gin_server

import "github.com/gin-contrib/cors"

type GinConfig struct {
	Addr string

	EnableBaseMw bool
	EnablePprof  bool
	EnableCROS   bool

	ProfilePath     string
	CORS            *cors.Config
	GracefulExitSec int64
}
