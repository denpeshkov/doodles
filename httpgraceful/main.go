package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server listen: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(shCtx); err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
