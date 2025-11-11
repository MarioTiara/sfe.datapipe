package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type LocalFolderLoader struct {
	Parser FileParser
	Path   string
}

func NewLocalFolderLoader(parser FileParser, path string) *LocalFolderLoader {
	return &LocalFolderLoader{
		Parser: parser,
		Path:   path,
	}
}

func (l *LocalFolderLoader) StreamRows() (<-chan []string, <-chan error) {
	dataCh := make(chan []string, 100)
	errCh := make(chan error, 10)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		files, err := os.ReadDir(l.Path)
		if err != nil {
			errCh <- err
			return
		}

		var wg sync.WaitGroup

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			wg.Add(1)
			go func(fileName string) {
				defer wg.Done()

				rowsCh, parserErrCh := l.Parser.Parse(filepath.Join(l.Path, fileName))

				for row := range rowsCh {
					dataCh <- row
				}

				for perr := range parserErrCh {
					if perr != nil {
						errCh <- fmt.Errorf("file %s: %w", fileName, perr)
					}
				}
			}(f.Name())
		}

		wg.Wait() // wait for all file parsers to finish
	}()

	return dataCh, errCh
}
