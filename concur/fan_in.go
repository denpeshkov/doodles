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
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						return
					}
					select {
					case mergedCh <- v:
					case <-ctx.Done():
						return
					}
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

// MergeTwo merges two input channels into one output channel. The order of the produced elements is undefined.
// The returned channel is closed when both input channels are closed or the provided context is canceled.
func MergeTwo[T any](ctx context.Context, ch1, ch2 <-chan T) <-chan T {
	mergedCh := make(chan T)

	go func() {
		defer close(mergedCh)
		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					break
				}
				mergedCh <- v
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					break
				}
				mergedCh <- v
			}
		}
	}()

	return mergedCh
}
