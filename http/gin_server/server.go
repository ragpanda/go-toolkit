package gin_server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type GinHttpServer struct {
	config *GinHttpConfig

	once   sync.Once
	server *http.Server
	engine *gin.Engine
}

func NewGinHttpServer(config *GinHttpConfig) *GinHttpServer {
	serv := GinHttpServer{}
	return &serv
}

type GinHttpConfig struct {
	Addr string

	EnableBaseMw bool
	EnablePprof  bool
	CORS         CORSCOnfig
}

type CORSCOnfig struct {
	Enable bool
	cors.Config
}

func (self *GinHttpServer) fillDefault(config *GinHttpConfig) {
	if config.Addr == "" {
		config.Addr = ":8080"
	}
}

func (self *GinHttpServer) Init() *GinHttpServer {
	self.once.Do(func() {
		self.fillDefault(self.config)
		self.engine = gin.Default()

		if self.config.EnablePprof {
			pprof.Register(self.engine)
		}
		if self.config.EnableBaseMw {
			self.engine.Use(BizDataMw, LogIDMW, StatMW)
		}

		if self.config.CORS.Enable {
			self.engine.Use(cors.New(self.config.CORS.Config))
		}

		self.server = &http.Server{
			Addr:           self.config.Addr,
			Handler:        self.engine,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	})
	return self
}

func (self *GinHttpServer) Run() error {
	return self.server.ListenAndServe()
}
