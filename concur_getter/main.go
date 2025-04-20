package main

import (
	"context"
	"errors"
	"sync"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Get1 calls [Getter.Get] for each address in parallel.
// Returns the first successful response. If all requests fail, returns an error.
func Get1(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	valCh := make(chan string, 1) // buffer of 1 to prevent missing the result
	errCh := make(chan error, 1)
	errCh <- nil // pass nil to prevent blocking on first

	var wg sync.WaitGroup
	wg.Add(len(addresses))
	for _, a := range addresses {
		go func() {
			defer wg.Done()

			if val, err := getter.Get(ctx, a, key); err != nil {
				errCh <- errors.Join(<-errCh, err)
			} else {
				select {
				case valCh <- val:
				default: // drop all but the first result to prevent blocking
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(valCh)
		close(errCh) // close to prevent blocking in case of empty addresses (no val and no err)
	}()

	select {
	case val, ok := <-valCh:
		if ok {
			return val, nil
		}
		err := <-errCh
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Get2 calls [Getter.Get] for each address in parallel.
// Returns the first successful response. If all requests fail, returns an error.
func Get2(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	if len(addresses) == 0 {
		return "", nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	valCh := make(chan string, 1) // buffer of 1 to prevent missing the result
	errCh := make(chan error, len(addresses))

	for _, a := range addresses {
		go func() {
			if val, err := getter.Get(ctx, a, key); err != nil {
				errCh <- err
			} else {
				select {
				case valCh <- val:
				default: // drop all but the first result to prevent blocking
				}
			}
		}()
	}

	var (
		joinErr error
		errCnt  int
	)
	for {
		select {
		case val := <-valCh:
			return val, nil
		case err := <-errCh:
			errCnt++
			joinErr = errors.Join(joinErr, err)
			if errCnt == len(addresses) {
				return "", joinErr
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
