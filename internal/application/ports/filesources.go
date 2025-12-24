package ports

import "io"

type FileInfo struct {
	Name string
	Size int64
}

// FileSource defines the abstraction for any file storage
type FileSource interface {
	// List files with a specific prefix
	List(prefix string) ([]FileInfo, error)

	// Open a file for reading (stream)
	Open(name string) (io.ReadCloser, error)

	// Move a file to a new destination
	Move(srcName string, destName string) error

	// Remove/delete a file
	Remove(name string) error
}
