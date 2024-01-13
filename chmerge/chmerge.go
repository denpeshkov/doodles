package chmerge

import "sync"

// Merge merges input channels into the output channel.
// The returned channel is closed when all input channels are closed.
func Merge[T any](channels ...<-chan T) <-chan T {
	ch := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, c := range channels {
		go func() {
			for v := range c {
				ch <- v
			}
			wg.Done()
		}()
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}

// FromSlice allows the use of a slice of bidirectional channels as an input to [Merge].
func FromSlice[T any](s []chan T) []<-chan T {
	ss := make([]<-chan T, 0, len(s))
	for _, v := range s {
		ss = append(ss, v)
	}
	return ss
}
