package service

import (
	"context"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ragpanda/go-toolkit/bizerr"
	"github.com/ragpanda/go-toolkit/log"
	"github.com/ragpanda/go-toolkit/utils"
)

type CommonResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	ErrCode string `json:"err_code"`
}
type HandlerFunc[Req any, Resp any] func(ctx context.Context, req *Req) (resp *Resp, err error)

func NewHandler[Req any, Resp any](path string) *Handler[Req, Resp] {
	return &Handler[Req, Resp]{
		path: path,
	}
}

type Handler[Req any, Resp any] struct {
	handler HandlerFunc[Req, Resp]
	path    string
}

func (self *Handler[Req, Resp]) WithHandler(handler HandlerFunc[Req, Resp]) *Handler[Req, Resp] {
	self.handler = handler
	return self
}

func (self Handler[Req, Resp]) Register(engine *gin.RouterGroup) gin.HandlerFunc {
	f := func(ctx *gin.Context) {
		req := new(Req)
		if ctx.Request.Method == "GET" {
			if err := ctx.BindQuery(&req); err != nil {
				ctx.JSON(400, &CommonResp{
					ErrCode: string(bizerr.ErrInvalidInput.Code()),
				})
				return
			}
		} else {
			if ctx == nil || ctx.Request == nil {
				log.Error(ctx, "ctx or ctx.Request is nil, %+v", ctx)
				ctx.JSON(400, &CommonResp{
					ErrCode: string(bizerr.ErrInvalidInput.Code()),
				})
				return
			}
			data, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Error(ctx, "read request body error: %v", err.Error())
				ctx.JSON(400, &CommonResp{
					ErrCode: string(bizerr.ErrInvalidInput.Code()),
				})
				return
			}
			err = utils.Unmarshal(data, &req)
			if err != nil {
				log.Error(ctx, "unmarshal request body error:`%v`, data:`%s`", err.Error(), data)
				ctx.JSON(400, &CommonResp{
					ErrCode: string(bizerr.ErrInvalidInput.Code()),
				})
				return
			}

			log.Debug(ctx, "request body: %s", data)
		}

		var resp *Resp
		err := utils.ProtectPanic(ctx, func() error {
			var err error
			resp, err = self.handler(ctx, req)
			return err
		})
		if err != nil {
			log.Warn(ctx, "handler error: %v", err.Error())
			var bizErr bizerr.BusinessError
			if errors.As(err, &bizErr) {
				ctx.JSON(500, &CommonResp{
					ErrCode: string(bizErr.Code()),
					Message: bizErr.Message(),
				})
				return
			}

			ctx.JSON(500, &CommonResp{
				ErrCode: string(bizerr.ErrInternalError.Code()),
				Message: err.Error(),
			})
			return
		}
		log.Debug(ctx, "handler success: %+v", resp)
		ctx.JSON(200, &CommonResp{
			Success: true,
			Message: "success",
			Data:    resp,
			ErrCode: "",
		})
		return
	}

	engine.POST(self.path, f)
	return f
}
