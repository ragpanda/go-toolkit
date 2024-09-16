package bizerr

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	biz2 "github.com/ragpanda/go-toolkit/biz"
)

// BusinessError 定义业务错误接口
type BusinessError interface {
	error
	Message() string
	Code() ErrorCode
	Is(error) bool
	WithMessage(message string) BusinessError
	WithStack(ctx context.Context) BusinessError
	StackTrace() string
}

// businessError 实现 BusinessError 接口
type businessError struct {
	code       ErrorCode
	message    string
	stackTrace string
	ctx        context.Context
}

// NewBusinessError 创建新的业务错误
func NewBusinessError(code ErrorCode, message string) BusinessError {
	return &businessError{
		code:    code,
		message: message,
	}
}

// Error 实现 error 接口
func (e *businessError) Error() string {
	var basic, ctxStr, stack string
	if e.stackTrace != "" {
		stack = fmt.Sprintf("stack:\n%s\n ", e.stackTrace)
	}

	if e.ctx != nil {
		bizData := biz2.GetBizData(e.ctx)
		ctxStr = fmt.Sprintf("ctx:%v ", bizData.String())
	}

	basic = fmt.Sprintf("code:%s, message:`%s` ", e.code, e.message)

	return fmt.Sprintf("{%s%s%s}", basic, ctxStr, stack)
}

func (e *businessError) Message() string {
	return e.message
}

// Code 返回错误码
func (e *businessError) Code() ErrorCode {
	return e.code
}

// Is 实现错误比较
func (e *businessError) Is(target error) bool {
	t, ok := target.(*businessError)
	if !ok {
		return false
	}
	return e.code == t.code
}

// WithMessage 创建一个新的错误，保持原始错误码但使用新的错误消息
func (e *businessError) WithMessage(message string) BusinessError {
	return &businessError{
		code:    e.code,
		message: message,
	}
}

// WithStack 附加调用栈信息
func (e *businessError) WithStack(ctx context.Context) BusinessError {
	if ctx != nil {
		e.ctx = ctx
	}

	stackTrace := getStackTrace()
	return &businessError{
		code:       e.code,
		message:    e.message,
		stackTrace: stackTrace,
	}
}

// StackTrace 返回调用栈信息
func (e *businessError) StackTrace() string {
	return e.stackTrace
}

// Is 检查两个错误是否相等
func Is(err, target error) bool {
	if err == nil {
		return target == nil
	}
	if target == nil {
		return false
	}
	be, ok := err.(BusinessError)
	if !ok {
		return false
	}
	return be.Is(target)
}

// getStackTrace 获取调用栈信息
func getStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	var builder strings.Builder
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return builder.String()
}
