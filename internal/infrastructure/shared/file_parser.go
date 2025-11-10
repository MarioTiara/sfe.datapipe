package shared

type FileParser interface {
	Parse(filePath string) (<-chan []string, <-chan error)
}
