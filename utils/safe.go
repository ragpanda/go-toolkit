package utils

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/ragpanda/go-toolkit/log"
)

func ProtectPanic(ctx context.Context, execFunc func() error) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			goErr := errors.Wrap(_err, 3)
			log.Error(ctx, "[ProtectPanic] panic recovered, err=%s\n%s", goErr.Error(), goErr.Stack())
			err = goErr
		}

		if err != nil {
			log.Debug(ctx, "[ProtectPanic] occur error %+v", err)
		}

	}()

	return execFunc()
}
