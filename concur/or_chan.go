package pattern

// Or returns a channel that is closed when any of the input channels is closed.
func Or(channels ...<-chan struct{}) <-chan struct{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	done := make(chan struct{})
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
