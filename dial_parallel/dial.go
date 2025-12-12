package main

import (
	"context"
	"errors"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, addr string) (net.Conn, error)
}

// Returns the first successful response. If all requests fail, returns an error.
func Dial(ctx context.Context, dialer Dialer, addrs []string) (net.Conn, error) {
	if len(addrs) == 0 {
		return nil, errors.New("empty addresses")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type dialResult struct {
		conn net.Conn
		err  error
	}
	dialResultCh := make(chan dialResult)

	for _, addr := range addrs {
		go func() {
			conn, err := dialer.Dial(ctx, addr)
			select {
			case dialResultCh <- dialResult{conn, err}:
			case <-ctx.Done():
				if conn != nil {
					_ = conn.Close()
				}
			}
		}()
	}

	var firstErr error
	for range len(addrs) {
		select {
		case res := <-dialResultCh:
			if res.conn != nil {
				return res.conn, nil
			}
			if firstErr == nil {
				firstErr = res.err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, firstErr
}
