package pattern

// ToChan returns a channel containing all elements in the slice s.
func ToChan[T any](s ...T) <-chan T {
	ch := make(chan T, len(s))
	for _, e := range s {
		ch <- e
	}
	close(ch)
	return ch
}
