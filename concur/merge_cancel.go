package pattern

import "context"

// MergeCancel returns a context that contains the values of ctx,
// and which is canceled when either ctx or cancelCtx is canceled.
func MergeCancel(ctx, cancelCtx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(ctx)
	stop := context.AfterFunc(cancelCtx, func() {
		cancel(context.Cause(cancelCtx))
	})
	return ctx, func() {
		stop()
		cancel(context.Canceled)
	}
}
