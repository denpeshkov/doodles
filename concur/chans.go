// Package concur provides common concurrency patterns.
package concur

import (
	"context"
	"sync"
)

// FanIn merges multiple input channels into one output channel. The order of the produced elements is undefined.
// The returned channel is closed when all input channels are closed or the provided context is canceled.
func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
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

// ToChan returns a channel containing all elements in the slice s.
func ToChan[T any](s ...T) <-chan T {
	ch := make(chan T, len(s))
	for _, e := range s {
		ch <- e
	}
	close(ch)
	return ch
}

// Or returns a channel that is closed when any of the input channels is closed.
func Or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	done := make(chan any)
	go func() {
		defer close(done)

		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-Or(append(channels[2:], done)...):
		}
	}()
	return done
}
