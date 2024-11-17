package utils

import (
	"context"
	"time"

	"github.com/go-errors/errors"
	"github.com/ragpanda/go-toolkit/log"
)

type RetryOption func(*RetryParams)

type RetryParams struct {
	RetryFailedInterval time.Duration
	MaxTimout           time.Duration
	MaxTimes            int
}

func RetryMaxTimout(retryMaxTimout, interval time.Duration) func(*RetryParams) {
	return func(params *RetryParams) {
		params.MaxTimout = retryMaxTimout
		params.RetryFailedInterval = interval
	}
}

func RetryMaxTimes(maxTimes int, interval time.Duration) func(*RetryParams) {
	return func(params *RetryParams) {
		params.MaxTimes = maxTimes
		params.RetryFailedInterval = interval
	}
}

func Retry(ctx context.Context, execFunc func(ctx context.Context) error, options ...RetryOption) error {
	p := &RetryParams{
		MaxTimout:           0,
		RetryFailedInterval: 1 * time.Second,
		MaxTimes:            3,
	}

	for _, op := range options {
		op(p)
	}

	if p.MaxTimes == 0 {
		p.MaxTimes = 1
	}

	var err error
	start := time.Now()

	var allowContinue func() bool
	if p.MaxTimout != 0 {
		allowContinue = func() bool {
			return time.Now().Sub(start) < p.MaxTimout
		}
	} else {
		t := 0
		allowContinue = func() bool {
			t++
			return t <= p.MaxTimes
		}
	}

	//仅统计
	execTimes := 0
	for allowContinue() {
		execTimes++
		err = panicSafe(ctx, func(ctx context.Context) error { return execFunc(ctx) })
		if err != nil {
			log.Info(ctx, "exec fail, times=%d, err=`%s`", execTimes, err.Error())
			time.Sleep(p.RetryFailedInterval)
			continue
		}
		return nil
	}

	return err
}

func panicSafe(ctx context.Context, execFunc func(ctx context.Context) error) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			goErr := errors.Wrap(_err, 3)
			log.Error(ctx, "[panicSafe] panic recovered, err=%s\n%s", goErr.Error(), goErr.Stack())
			err = goErr
		}
	}()

	return execFunc(ctx)
}
