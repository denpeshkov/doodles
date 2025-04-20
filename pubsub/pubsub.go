// Package pubsub implements a Pub/Sub.
package pubsub

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"slices"
)

var errPubSubClosed = errors.New("PubSub closed")

// Subscription is a [PubSub] subscription.
type Subscription[T any] struct {
	ch chan T
}

// Updates returns a channel to receive messages from [PubSub].
func (s *Subscription[T]) Updates() <-chan T { return s.ch }

// PubSub is a Pub/Sub system.
type PubSub[T any] struct {
	subs     []Subscription[T]
	actch    chan func()
	closedch chan struct{}
}

// NewPubSub creates and returns a new [PubSub] instance.
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		actch:    make(chan func()),
		closedch: make(chan struct{}),
	}
}

// Run starts the [PubSub] message processing loop.
// It blocks until the provided context is canceled, at which point
// it cleans up resources and closes all subscriptions.
// This method must be called before any other methods on the [PubSub].
func (ps *PubSub[T]) Run(ctx context.Context) {
	defer close(ps.closedch)
	for {
		select {
		case f := <-ps.actch:
			f()
		case <-ctx.Done():
			for _, sub := range ps.subs {
				close(sub.ch)
			}
			clear(ps.subs)
			return
		}
	}
}

// Subscribe creates and returns a new subscription with the specified buffer size.
func (ps *PubSub[T]) Subscribe(bufSize int) (Subscription[T], error) {
	sub := Subscription[T]{ch: make(chan T, bufSize)}
	if err := ps.process(func() {
		ps.subs = append(ps.subs, sub)
	}); err != nil {
		return Subscription[T]{}, err
	}
	return sub, nil
}

// Unsubscribe removes the given subscription from the [PubSub].
func (ps *PubSub[T]) Unsubscribe(sub Subscription[T]) {
	_ = ps.process(func() {
		// Can be optimized by swapping the element to be deleted
		// with the last element, as we don't need to preserve the order.
		ps.subs = slices.DeleteFunc(ps.subs, func(s Subscription[T]) bool { return s == sub })
		close(sub.ch)
	})
}

// Publish publishes the given message to all active subscriptions.
func (ps *PubSub[T]) Publish(ctx context.Context, msg T) error {
	ch := make(chan error)
	if err := ps.process(func() {
		var err error
	loop:
		for _, sub := range ps.subs {
			select {
			case sub.ch <- msg:
			case <-ctx.Done():
				err = errors.Join(err, fmt.Errorf("message undelivered: %w", ctx.Err()))
				break loop
			}
		}
		ch <- err
	}); err != nil {
		return err
	}
	return <-ch
}

// Subscriptions returns an iterator over the current subscriptions.
func (ps *PubSub[T]) Subscriptions() iter.Seq[Subscription[T]] {
	// Make a snapshot.
	ch := make(chan []Subscription[T])
	if err := ps.process(func() {
		ch <- slices.Clone(ps.subs)
	}); err != nil {
		return nil
	}
	subs := <-ch

	return func(yield func(Subscription[T]) bool) {
		for _, sub := range subs {
			if !yield(sub) {
				return
			}
		}
	}
}

func (ps *PubSub[T]) process(f func()) error {
	select {
	case ps.actch <- f:
	case <-ps.closedch:
		return errPubSubClosed
	}
	return nil
}
