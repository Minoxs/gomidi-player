package utils

import (
	"context"
	"errors"
)

var ErrUnavailable = errors.New("unavailable context")

// UnavailableContext returns a context that is already canceled
func UnavailableContext() context.Context {
	// Create a context and cancel immediately
	var ctx, cancel = context.WithCancelCause(context.TODO())
	cancel(ErrUnavailable)
	return ctx
}
