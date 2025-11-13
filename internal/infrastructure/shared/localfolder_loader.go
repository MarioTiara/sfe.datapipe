package shared

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type LocalFolderLoader struct {
	logger logger.Logger
	Parser FileParser
	Path   string
}

func NewLocalFolderLoader(logger logger.Logger, parser FileParser, path string) *LocalFolderLoader {
	return &LocalFolderLoader{
		logger: logger,
		Parser: parser,
		Path:   path,
	}
}

func (l *LocalFolderLoader) StreamRows(ctx context.Context) (<-chan []string, <-chan error) {
	dataCh := make(chan []string, 100)
	errCh := make(chan error, 10)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		// 🔹 Log start of the loading process
		l.logger.Info(ctx, "Starting to stream files from folder",
			logger.Field{Key: "path", Value: l.Path})

		files, err := os.ReadDir(l.Path)
		if err != nil {
			l.logger.Error(ctx, "Failed to read directory", err,
				logger.Field{Key: "path", Value: l.Path})
			errCh <- err
			return
		}

		if len(files) == 0 {
			l.logger.Info(ctx, "No files found in folder",
				logger.Field{Key: "path", Value: l.Path})
		}

		var wg sync.WaitGroup

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			filePath := filepath.Join(l.Path, f.Name())
			wg.Add(1)

			go func(fileName, fullPath string) {
				defer wg.Done()

				// 🔹 Log before parsing
				l.logger.Info(ctx, "Parsing file started",
					logger.Field{Key: "file", Value: fileName})

				rowsCh, parserErrCh := l.Parser.Parse(fullPath)

				for row := range rowsCh {
					select {
					case <-ctx.Done():
						l.logger.Error(ctx, "Context cancelled while reading rows", ctx.Err(),
							logger.Field{Key: "file", Value: fileName})
						return
					case dataCh <- row:
					}
				}

				for perr := range parserErrCh {
					if perr != nil {
						l.logger.Error(ctx, "Parser error occurred", perr,
							logger.Field{Key: "file", Value: fileName})
						errCh <- fmt.Errorf("file %s: %w", fileName, perr)
					}
				}

				// 🔹 Log after parsing
				l.logger.Info(ctx, "Parsing file completed",
					logger.Field{Key: "file", Value: fileName})

			}(f.Name(), filePath)
		}

		wg.Wait()

		// 🔹 Log after all files processed
		l.logger.Info(ctx, "Finished streaming all files from folder",
			logger.Field{Key: "path", Value: l.Path})
	}()

	return dataCh, errCh
}
