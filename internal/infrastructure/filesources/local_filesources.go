package ports

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ports"
)

type LocalFileSource struct {
	BaseDir string
}

func (l *LocalFileSource) List(prefix string) ([]ports.FileInfo, error) {
	entries, err := os.ReadDir(l.BaseDir)
	if err != nil {
		return nil, err
	}

	var result []ports.FileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if strings.Contains(e.Name(), prefix) {
			info, _ := e.Info()
			result = append(result, ports.FileInfo{
				Name: e.Name(),
				Size: info.Size(),
			})
		}
	}

	return result, nil
}

func (l *LocalFileSource) Open(name string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.BaseDir, name))
}

func (l *LocalFileSource) Move(srcName, destFolder string) error {
	srcPath := filepath.Join(l.BaseDir, srcName)

	fileName := filepath.Base(srcName) // keep same filename
	destPath := filepath.Join(l.BaseDir, destFolder, fileName)

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	return os.Rename(srcPath, destPath)
}

// Remove/delete a file
func (l *LocalFileSource) Remove(name string) error {
	path := filepath.Join(l.BaseDir, name)
	return os.Remove(path)
}
