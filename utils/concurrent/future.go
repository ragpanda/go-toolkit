package concurrent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ragpanda/go-toolkit/log"
	"github.com/ragpanda/go-toolkit/utils"
)

type Future[T any] struct {
	once sync.Once

	done     chan struct{}
	lock     sync.RWMutex
	result   T
	err      error
	execFunc func(ctx context.Context) (T, error)
}

func NewFuture[T any]() *Future[T] {
	f := &Future[T]{}
	f.done = make(chan struct{}, 1)
	return f
}

func (f *Future[T]) SetResult(result T) {
	f.setResultOrError(&result, nil)
}

func (f *Future[T]) SetError(err error) {
	f.setResultOrError(nil, err)
}

func (f *Future[T]) IsDone() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

func (f *Future[T]) Error() error {
	<-f.done
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.err
}

func (f *Future[T]) Result() T {
	<-f.done
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.result
}

func (f *Future[T]) Get() (T, error) {
	<-f.done
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.result, f.err
}

func (f *Future[T]) GetResultWithTimeout(timeout time.Duration) (T, error) {
	select {
	case <-f.done:
		return f.Get()
	case <-time.After(timeout):
		return f.result, fmt.Errorf("wait result timeout")
	}
}

func (self *Future[T]) SetExec(exec func(ctx context.Context) (T, error)) *Future[T] {
	self.execFunc = exec
	return self
}

func (self *Future[T]) Exec(ctx context.Context) *Future[T] {
	err := utils.ProtectPanic(ctx, func() error {
		result, err := self.execFunc(ctx)
		if err != nil {
			return err
		}
		self.SetResult(result)
		return nil
	})
	if err != nil {
		self.SetError(err)
	}
	return self
}

func (f *Future[T]) setResultOrError(result *T, err error) {
	f.once.Do(func() {
		f.lock.Lock()
		defer f.lock.Unlock()
		if err != nil {
			f.err = err
		} else {
			f.result = *result
		}
	})

	select {
	case f.done <- struct{}{}:
		close(f.done)
	default:
		panic("Future already done")
	}
}

func RunWithFuture[T any](ctx context.Context, exec func(ctx context.Context) (T, error)) *Future[T] {
	future := NewFuture[T]()
	future.SetExec(exec)
	go func() error {
		future.Exec(ctx)
		return nil
	}()

	return future
}

type AsyncPool[T any] struct {
	ctx       context.Context
	workerNum int
	once      sync.Once
	taskChan  chan *Future[T]
	cancel    context.CancelFunc
	waitGroup sync.WaitGroup
}

func NewAsyncPool[T any](ctx context.Context, workerNum int) *AsyncPool[T] {
	if workerNum <= 0 {
		workerNum = 1
	}

	p := &AsyncPool[T]{
		workerNum: workerNum,
	}
	p.taskChan = make(chan *Future[T], workerNum*10)
	p.ctx, p.cancel = context.WithCancel(ctx)
	return p
}

func (pool *AsyncPool[T]) Go(ctx context.Context, exec func(ctx context.Context) (T, error)) *Future[T] {
	pool.once.Do(func() {
		for i := 0; i < pool.workerNum; i++ {
			pool.waitGroup.Add(1)
			go func(number int) {
				defer pool.waitGroup.Done()
				log.Debug(ctx, "async pool run, gonum:%d", number)
				for {
					select {
					case task := <-pool.taskChan:
						log.Debug(ctx, "async pool exec, gonum:%d", number)
						task.Exec(ctx)
						log.Debug(ctx, "async pool done, gonum:%d", number)
					case <-pool.ctx.Done():
						log.Debug(ctx, "async pool exit, gonum:%d", number)
						return
					}
				}
			}(i)
		}
	})
	future := NewFuture[T]()
	future.SetExec(exec)
	pool.taskChan <- future
	return future
}

func (pool *AsyncPool[T]) Close() {
	pool.cancel()
	pool.waitGroup.Wait()
}
