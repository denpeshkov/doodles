package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"testing"
	"time"

	"go.uber.org/goleak"
)

type resp struct {
	val   string
	err   error
	delay time.Duration
}

type getter struct {
	responses map[string]map[string]resp
}

func (m *getter) Get(ctx context.Context, addr, key string) (string, error) {
	resp, ok := m.responses[addr][key]
	if !ok {
		return "", fmt.Errorf("address %q key %q is missing", addr, key)
	}
	if resp.delay > 0 {
		select {
		case <-time.After(resp.delay):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return resp.val, resp.err
}

func TestGet(t *testing.T) {
	t.Run("Get1", func(t *testing.T) {
		testGet(t, Get1)
	})
	t.Run("Get2", func(t *testing.T) {
		testGet(t, Get2)
	})
}

func TestGet_Error(t *testing.T) {
	t.Run("Get1", func(t *testing.T) {
		testGet_Error(t, Get1)
	})
	t.Run("Get2", func(t *testing.T) {
		testGet_Error(t, Get2)
	})
}

func testGet(t *testing.T, get func(ctx context.Context, getter Getter, addresses []string, key string) (string, error)) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name      string
		resps     map[string]map[string]resp
		key       string
		ttl       time.Duration
		wantValue string
	}{
		{
			name: "first address fails second succeeds",
			resps: map[string]map[string]resp{
				"addr1": {
					"key1": {err: errors.New("connection error")},
				},
				"addr2": {
					"key1": {val: "value2"},
				},
			},
			key:       "key1",
			wantValue: "value2",
			ttl:       1 * time.Millisecond,
		},
		{
			name: "fast address wins over slow",
			resps: map[string]map[string]resp{
				"addr1": {
					"key1": {val: "value1", delay: 200 * time.Millisecond},
				},
				"addr2": {
					"key1": {val: "value2", delay: 50 * time.Millisecond},
				},
			},
			key:       "key1",
			ttl:       100 * time.Millisecond,
			wantValue: "value2",
		},
		{
			name:      "empty address list",
			resps:     map[string]map[string]resp{},
			key:       "key1",
			ttl:       50 * time.Millisecond,
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := &getter{tt.resps}

			ctx, cancel := context.WithTimeout(context.Background(), tt.ttl)
			t.Cleanup(cancel)

			addresses := slices.Collect(maps.Keys(tt.resps))
			got, err := get(ctx, getter, addresses, tt.key)
			if err != nil {
				t.Fatalf("Get() failed: %v", err)
			}
			if got != tt.wantValue {
				t.Errorf("Get() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func testGet_Error(t *testing.T, get func(ctx context.Context, getter Getter, addresses []string, key string) (string, error)) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name    string
		resps   map[string]map[string]resp
		key     string
		ttl     time.Duration
		testErr func(t *testing.T, err error)
	}{

		{
			name: "all addresses fail",
			resps: map[string]map[string]resp{
				"addr1": {
					"key1": {err: errors.New("error 1")},
				},
				"addr2": {
					"key1": {err: errors.New("error 2")},
				},
			},
			key:     "key1",
			ttl:     1 * time.Millisecond,
			testErr: testGetErrors,
		},
		{
			name: "context cancellation",
			resps: map[string]map[string]resp{
				"addr1": {
					"key1": {val: "value1", delay: 200 * time.Millisecond},
				},
			},
			key:     "key1",
			ttl:     50 * time.Millisecond,
			testErr: testCtxDeadlineErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := &getter{tt.resps}

			ctx, cancel := context.WithTimeout(context.Background(), tt.ttl)
			t.Cleanup(cancel)

			addresses := slices.Collect(maps.Keys(tt.resps))
			_, err := get(ctx, getter, addresses, tt.key)
			tt.testErr(t, err)
		})
	}
}

func testGetErrors(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected Get() to return an error")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Get() returned unexpected error: %v", context.DeadlineExceeded)
	}
}

func testCtxDeadlineErr(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Get() = %v, want %v", err, context.DeadlineExceeded)
	}
}
