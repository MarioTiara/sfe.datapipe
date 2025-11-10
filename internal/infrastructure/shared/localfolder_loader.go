package shared

import (
	"fmt"
	"os"
	"path/filepath"
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
	dataCh := make(chan []string)
	errCh := make(chan error, 1)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		files, err := os.ReadDir(l.Path)
		mess := fmt.Sprintf("files found: %d", len(files))
		fmt.Println(mess)
		if err != nil {
			errCh <- err
			return
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			fullPath := filepath.Join(l.Path, f.Name())
			rowsCh, parserErrCh := l.Parser.Parse(fullPath)
			mess := fmt.Sprintf("files found: %s", fullPath)
			fmt.Println(mess)
			for row := range rowsCh {
				dataCh <- row
			}

			if err := <-parserErrCh; err != nil {
				errCh <- fmt.Errorf("file %s: %w", f.Name(), err)
				return
			}
		}
	}()

	fmt.Println("local data loader finished")
	return dataCh, errCh
}
