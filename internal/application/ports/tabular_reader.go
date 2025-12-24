package ports

type TabularFileReader interface {
	ReadRows(name string) (<-chan []string, error)
}
