package concurrent

import (
	"context"
	"sync"
	"time"

	"github.com/ragpanda/go-toolkit/utils"
)

// merge multi channel to one channel, []chan T -> chan T
func MergeChannels[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	output := make(chan T)

	// Launch a goroutine for each input channel
	for i, _c := range channels {
		wg.Add(1)
		go func(channel <-chan T, idx int) {
			defer wg.Done()
			for value := range channel {
				output <- value
			}
		}(_c, i)
	}

	// Close the output channel once all input channels have been processed
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// MergeChannelsWithBatch merge channel and merge data to batch,  [] chan T -> chan []T
// flushSize: max batch size, if reach this size, will flush immediately
// flushInterval: max flush interval, if reach this interval, will flush immediately
// Note: the sort of batch is not guaranteed
func MergeChannelsWithBatch[T any](
	ctx context.Context,
	flushSize int, flushInterval time.Duration,
	channels ...<-chan T,
) <-chan []T {
	out := MergeChannels(ctx, channels...)
	outputWithBuffer := make(chan []T)
	go utils.ProtectPanic(ctx, func() error {
		defer close(outputWithBuffer)

		run := true
		flush := false

		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()
		buffered := make([]T, 0, flushSize+1)
		for run {
			select {
			case <-ctx.Done():
				run = false
				flush = true
			case m, hasMore := <-out:
				if !hasMore {
					run = false
					flush = true
					break
				}

				buffered = append(buffered, m)
				if len(buffered) >= flushSize {
					flush = true
				}
			case <-ticker.C:
				flush = true
			}

			if len(buffered) != 0 && flush {
				ticker.Reset(flushInterval)
				outputWithBuffer <- buffered
				buffered = make([]T, 0, flushSize+1)
			}
		}

		return nil
	})
	return outputWithBuffer
}
