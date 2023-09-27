package concurrent

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMergeChannelsWithBatch(t *testing.T) {
	testCases := []struct {
		channels        []chan int    // input channels
		flushSize       int           // number of elements to flush
		flushInterval   time.Duration // duration between flushes
		expectedBatches []int         // expected output batches
	}{
		{
			channels: []chan int{
				make(chan int, 2),
				make(chan int, 2),
				make(chan int, 2),
			},
			flushSize:       2,
			flushInterval:   10 * time.Second,
			expectedBatches: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tc := range testCases {
		ctx, cancel := context.WithCancel(context.Background())

		// Populate the input channels with some test data
		channelsReadOnly := make([]<-chan int, 0, len(tc.channels))
		for i, _ch := range tc.channels {
			c := _ch
			go func(c chan int, offset int) {
				defer close(c)

				for j := 0; j < 2; j++ {
					c <- (offset + j + 1)
				}

			}(c, i*2) // 0, 2, 4
			channelsReadOnly = append(channelsReadOnly, c)
		}

		// Call the function under test
		out := MergeChannelsWithBatch(ctx, tc.flushSize, tc.flushInterval, channelsReadOnly...)

		// Verify that the output matches the expected batches
		elemList := make([]int, 0)
		for batch := range out {
			elemList = append(elemList, batch...)

		}
		sort.Ints(elemList)
		require.Equal(t, tc.expectedBatches, elemList)

		// Verify that the output channel is closed after all batches have been received
		_, ok := <-out
		require.False(t, ok)

		cancel()
	}
}
