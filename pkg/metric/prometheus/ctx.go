package prometheus

import (
	"context"
	"time"
)

type dataKey struct{}

type Data struct {
	Path    string
	Method  string
	StartAt time.Time
}

// NewContext creates a new context with the given Data.
func NewContext(ctx context.Context, data Data) context.Context {
	return context.WithValue(ctx, dataKey{}, data)
}

// FromContext returns Data from the given context.
func FromContext(ctx context.Context) (Data, bool) {
	c, ok := ctx.Value(dataKey{}).(Data)
	return c, ok
}
