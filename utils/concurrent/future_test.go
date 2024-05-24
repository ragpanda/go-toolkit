package concurrent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFutureGetResultReturnsResult(t *testing.T) {
	future := NewFuture[string]()
	go func() {
		future.SetResult("Hello, World!")
	}()

	result, err := future.Get()
	require.NoError(t, err)
	require.Equal(t, "Hello, World!", result)
}

func TestFutureGetResultReturnsError(t *testing.T) {
	future := NewFuture[string]()
	go func() {
		future.SetError(errors.New("an error occurred"))
	}()

	_, err := future.Get()
	require.Error(t, err)
	require.Equal(t, "an error occurred", err.Error())
}

func TestFutureGetResultWithTimeoutReturnsResultBeforeTimeout(t *testing.T) {
	future := NewFuture[string]()
	go func() {
		future.SetResult("Hello, World!")
	}()

	result, err := future.GetResultWithTimeout(2 * time.Second)
	require.NoError(t, err)
	require.Equal(t, "Hello, World!", result)
}

func TestFutureGetResultWithTimeoutReturnsErrorBeforeTimeout(t *testing.T) {
	future := NewFuture[string]()
	go func() {
		future.SetError(errors.New("an error occurred"))
	}()

	_, err := future.GetResultWithTimeout(2 * time.Second)
	require.Error(t, err)
	require.Equal(t, "an error occurred", err.Error())
}

func TestFutureGetResultWithTimeoutReturnsErrorAfterTimeout(t *testing.T) {
	future := NewFuture[string]()
	go func() {
		time.Sleep(3 * time.Second)
		future.SetResult("Hello, World!")
	}()

	_, err := future.GetResultWithTimeout(2 * time.Second)
	require.Error(t, err)
	require.Equal(t, "wait result timeout", err.Error())
}

func TestFutureExecReturnsResult(t *testing.T) {
	future := NewFuture[string]()
	future.SetExec(func(ctx context.Context) (string, error) {
		return "Hello, World!", nil
	})
	future.Exec(context.Background())

	result, err := future.Get()
	require.NoError(t, err)
	require.Equal(t, "Hello, World!", result)
}

func TestRunWithFutureExecutesSuccessfully(t *testing.T) {
	ctx := context.Background()
	future := RunWithFuture(ctx, func(ctx context.Context) (string, error) {
		return "Hello, World!", nil
	})

	result, err := future.Get()
	require.NoError(t, err)
	require.Equal(t, "Hello, World!", result)
}

func TestRunWithFutureHandlesError(t *testing.T) {
	ctx := context.Background()
	future := RunWithFuture(ctx, func(ctx context.Context) (string, error) {
		return "", errors.New("an error occurred")
	})

	_, err := future.Get()
	require.Error(t, err)
	require.Equal(t, "an error occurred", err.Error())
}

func TestRunWithFutureHandlesPanic(t *testing.T) {
	ctx := context.Background()
	future := RunWithFuture(ctx, func(ctx context.Context) (string, error) {
		panic("unexpected error")
	})

	_, err := future.Get()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected error")
}

func TestAsyncPoolExecutesTasks(t *testing.T) {
	ctx := context.Background()
	pool := NewAsyncPool[string](ctx, 2)

	var results []string
	futureList := make([]*Future[string], 0)
	for i := 0; i < 5; i++ {
		future := pool.Go(ctx, func(ctx context.Context) (string, error) {
			return "Hello, World!", nil
		})
		futureList = append(futureList, future)
	}

	for _, future := range futureList {
		result, err := future.Get()
		require.NoError(t, err)
		require.Equal(t, "Hello, World!", result)
		results = append(results, result)
	}

	require.Len(t, results, 5)

	pool.Close()
}
