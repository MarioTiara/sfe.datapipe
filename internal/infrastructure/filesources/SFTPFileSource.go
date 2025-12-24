package ports

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ports"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPFileSource struct {
	Client  *sftp.Client
	BaseDir string
}

func NewSFTPFileSource(user, password, host string, port int, baseDir string) (*SFTPFileSource, error) {
	// SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Connect to SSH
	addr := fmt.Sprintf("%s:%d", host, port)
	sshConn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect ssh: %w", err)
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create sftp client: %w", err)
	}

	return &SFTPFileSource{
		Client:  sftpClient,
		BaseDir: baseDir,
	}, nil
}

// List files with a given prefix
func (s *SFTPFileSource) List(prefix string) ([]ports.FileInfo, error) {
	// s.Client.Stat()
	entries, err := s.Client.ReadDir(s.BaseDir)
	if err != nil {
		return nil, err
	}

	var result []ports.FileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.Contains(e.Name(), prefix) {
			result = append(result, ports.FileInfo{
				Name: e.Name(),
				Size: e.Size(),
			})
		}
	}

	return result, nil
}

// Open a file for reading
func (s *SFTPFileSource) Open(name string) (io.ReadCloser, error) {
	path := path.Join(s.BaseDir, name)
	return s.Client.Open(path)
}

// Move a file to another folder, keeping the filename
func (s *SFTPFileSource) Move(srcName, destFolder string) error {
	srcPath := path.Join(s.BaseDir, srcName)
	fileName := path.Base(srcName)
	destPath := path.Join(s.BaseDir, destFolder, fileName)

	// Ensure destination directory exists
	if err := s.Client.MkdirAll(path.Join(s.BaseDir, destFolder)); err != nil {
		return err
	}

	return s.Client.Rename(srcPath, destPath)
}

// Remove/delete a file
func (s *SFTPFileSource) Remove(name string) error {
	path := path.Join(s.BaseDir, name)
	return s.Client.Remove(path)
}
