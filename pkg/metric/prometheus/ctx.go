package prometheus

import (
	"context"
	"net/http"
)

type dataKey struct{}

type Data struct {
	Path string
}

// NewContext creates a new context with the given Data.
func NewContext(ctx context.Context, data Data) context.Context {
	return context.WithValue(ctx, &dataKey{}, data)
}

// FromContext returns Data from the given context.
func FromContext(ctx context.Context) (Data, bool) {
	c, ok := ctx.Value(&dataKey{}).(Data)
	return c, ok
}

func NewHTTPRequestContext(r *http.Request, data Data) {
	ctx := context.WithValue(r.Context(), &dataKey{}, data)

	r2 := r.WithContext(ctx)

	*r = *r2
}
