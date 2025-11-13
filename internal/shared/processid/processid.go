package processid

import (
	"context"

	"github.com/google/uuid"
)

type key string

const processIDKey key = "process_id"

func New() string {
	return uuid.New().String()
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, processIDKey, New())
}

func FromContext(ctx context.Context) string {
	if v := ctx.Value(processIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}
