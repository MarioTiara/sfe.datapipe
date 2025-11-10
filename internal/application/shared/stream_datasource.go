package shared

type StreamDataSource interface {
	StreamRows() (<-chan []string, <-chan error)
}
