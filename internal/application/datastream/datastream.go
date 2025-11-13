package datastream

import "context"

type DataStream interface {
	StreamRows(ctx context.Context) (<-chan []string, <-chan error)
}
