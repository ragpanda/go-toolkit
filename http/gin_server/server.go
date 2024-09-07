package gin_server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ragpanda/go-toolkit/log"
)

type GinHttpServer struct {
	config *GinConfig

	once   sync.Once
	server *http.Server
	engine *gin.Engine
}

func NewGinHttpServer(config *GinConfig) *GinHttpServer {
	serv := GinHttpServer{}
	return &serv
}

func (self *GinHttpServer) Init() *GinHttpServer {
	self.once.Do(func() {
		self.fillDefault(self.config)
		self.engine = gin.Default()

		if self.config.EnablePprof {
			pprof.Register(self.engine, self.config.ProfilePath)
		}
		if self.config.EnableBaseMw {
			self.engine.Use(BizDataMw, StatMW)
		}

		if self.config.EnableCROS {
			self.engine.Use(cors.New(*self.config.CORS))
		}

		self.server = &http.Server{
			Addr:           self.config.Addr,
			Handler:        self.engine,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	})
	return self
}

func (self *GinHttpServer) Run(ctx context.Context) error {
	var e error
	positiveExit := make(chan struct{}, 1)
	go func() {
		defer func() {
			positiveExit <- struct{}{}
		}()
		err := self.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error(ctx, "http server start failed %s", err.Error())
			e = err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Info(ctx, "Server is shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(self.config.GracefulExitSec)*time.Second)
		defer cancel()
		if err := self.server.Shutdown(ctx); err != nil {
			log.Error(ctx, "Server Shutdown:", err)
			e = err
		}
	case <-positiveExit:
		log.Info(ctx, "Server is shutting down")

	}

	return e
}

func (self *GinHttpServer) GetEngine(group string) *gin.Engine {
	return self.engine
}

func (self *GinHttpServer) GetConfig() GinConfig {
	return *self.config
}

func (self *GinHttpServer) fillDefault(config *GinConfig) {
	if config.Addr == "" {
		config.Addr = ":8080"
	}
	if config.GracefulExitSec == 0 {
		config.GracefulExitSec = 10
	}
	if config.ProfilePath == "" {
		config.ProfilePath = "/debug/pprof"
	}
}
