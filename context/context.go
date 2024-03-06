package context

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type emptyContext struct{}

func (c *emptyContext) Deadline() (deadline time.Time, ok bool) { return }

func (c *emptyContext) Done() <-chan struct{} { return nil }

func (c *emptyContext) Err() error { return nil }

func (c *emptyContext) Value(key any) any { return nil }

type bgContext struct {
	emptyContext
}

func (c *bgContext) String() string { return "Background" }

func Background() Context {
	return &bgContext{}
}

type todoContext struct {
	emptyContext
}

func (c *todoContext) String() string { return "TODO" }

func TODO() Context {
	return &todoContext{}
}

type CancelFunc func()

var Canceled = errors.New("context canceled")

type cancelContext struct {
	Context
	done chan struct{}
	err  error
	mu   sync.Mutex
}

func (c *cancelContext) Done() <-chan struct{} { return c.done }

func (c *cancelContext) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *cancelContext) cancel(err error) {
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		close(c.done)
	}
}

func WithCancel(parent Context) (Context, CancelFunc) {
	ctx := &cancelContext{
		Context: parent,
		done:    make(chan struct{}),
		mu:      sync.Mutex{},
	}
	cancel := func() { ctx.cancel(Canceled) }

	go func() {
		select {
		case <-parent.Done():
			ctx.cancel(parent.Err())
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

var DeadlineExceeded = deadlineExceededErr{}

type deadlineExceededErr struct{}

func (deadlineExceededErr) Error() string   { return "deadline exceeded " }
func (deadlineExceededErr) Timeout() bool   { return true }
func (deadlineExceededErr) Temporary() bool { return true }

type deadlineContext struct {
	*cancelContext
	deadline time.Time
}

func (c *deadlineContext) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	cctx, cancel := WithCancel(parent)
	ctx := &deadlineContext{
		cancelContext: cctx.(*cancelContext),
		deadline:      d,
	}

	t := time.AfterFunc(time.Until(d), func() { ctx.cancel(DeadlineExceeded) })
	stop := func() {
		t.Stop()
		cancel()
	}
	return ctx, stop
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

type valueContext struct {
	Context
	key, value any
}

func (c *valueContext) Value(key any) any {
	if c.key == key {
		return c.value
	}
	return c.Context.Value(key)
}

func WithValue(parent Context, key, val any) Context {
	if key == nil {
		panic("key is nil")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	return &valueContext{
		Context: parent,
		key:     key,
		value:   val,
	}
}
