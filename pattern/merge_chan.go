package pattern

import (
	"context"
	"sync"
)

// Merge merges multiple input channels into one output channel. The order of the produced elements is undefined.
// The returned channel is closed when all input channels are closed or the provided context is canceled.
func Merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	mergedCh := make(chan T)
	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				select {
				case mergedCh <- v:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(mergedCh)
	}()

	return mergedCh
}
