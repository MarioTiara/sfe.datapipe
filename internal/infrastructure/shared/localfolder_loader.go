package shared

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type LocalFolderLoader struct {
	logger     logger.Logger
	Parser     FileParser
	Path       string
	ArcivePath string
}

func NewLocalFolderLoader(logger logger.Logger, parser FileParser, path string, archievePath string) *LocalFolderLoader {
	return &LocalFolderLoader{
		logger:     logger,
		Parser:     parser,
		Path:       path,
		ArcivePath: archievePath,
	}
}

func (l *LocalFolderLoader) StreamRows(ctx context.Context) (<-chan []string, <-chan error) {
	dataCh := make(chan []string, 100)
	errCh := make(chan error, 10)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		folderRead_start := time.Now()

		files, err := os.ReadDir(l.Path)
		if err != nil {
			l.logger.Error(ctx, "Failed to read directory", err,
				logger.Field{Key: "path", Value: l.Path})
			errCh <- err
			return
		}
		// 🔹 Log start of the loading process
		l.logger.Info(ctx, "Starting to stream files from folder",
			logger.Field{
				Key: "path", Value: l.Path,
			},
			logger.Field{
				Key: "files_count", Value: len(files),
			})

		if len(files) == 0 {
			l.logger.Info(ctx, "No files found in folder",
				logger.Field{Key: "path", Value: l.Path})
		}

		var wg sync.WaitGroup

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			file_start := time.Now()
			filePath := filepath.Join(l.Path, f.Name())
			wg.Add(1)

			go func(fileName, fullPath string) {
				defer wg.Done()
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

				l.logger.Info(ctx, "Parsing file completed",
					logger.Field{Key: "file", Value: fileName},
					logger.Field{Key: "processing_time", Value: time.Since(file_start).String()},
				)

				// 🔹 Move file to archive after parsing completes
				l.archieveData(ctx, fileName)

			}(f.Name(), filePath)
		}

		wg.Wait()
		// 🔹 Log after all files processed
		l.logger.Info(ctx, "Finished streaming all files from folder",
			logger.Field{Key: "path", Value: l.Path},
			logger.Field{Key: "total_time", Value: time.Since(folderRead_start).String()})
	}()

	return dataCh, errCh
}

func (l *LocalFolderLoader) archieveData(ctx context.Context, fileName string) {
	archiveFolder := l.ArcivePath
	if err := os.MkdirAll(archiveFolder, 0755); err != nil {
		l.logger.Error(ctx, "Error creating archive folder", err)
		return
	}

	oldPath := filepath.Join(l.Path, fileName)
	newFileName := l.createPrefix() + fileName
	newPath := filepath.Join(archiveFolder, newFileName)

	// Try to move file
	if err := os.Rename(oldPath, newPath); err != nil {
		l.logger.Error(ctx, "Error moving file to archive", err,
			logger.Field{Key: "file", Value: fileName})

		// Fallback: copy + remove (cross-drive safe)
		data, readErr := os.ReadFile(oldPath)
		if readErr == nil {
			if writeErr := os.WriteFile(newPath, data, 0644); writeErr == nil {
				_ = os.Remove(oldPath)
			} else {
				l.logger.Error(ctx, "Failed to write archive file", writeErr,
					logger.Field{Key: "file", Value: fileName})
			}
		}
	}
}

func (l *LocalFolderLoader) createPrefix() string {
	now := time.Now()

	formatted := now.Format("020106150405")

	return "Processed_" + formatted + "_"
}
