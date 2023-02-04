package storage

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/casdoor/oss"
)

var baseFolder = "files"

// FileSystem file system storage
type FileSystem struct {
	Base string
}

// NewFileSystem initialize the local file system storage
func NewFileSystem(base string) *FileSystem {
	absBase, err := filepath.Abs(base)
	if err != nil {
		panic("local file system storage's base folder is not initialized")
	}

	return &FileSystem{Base: absBase}
}

// GetFullPath get full path from absolute/relative path
func (fileSystem FileSystem) GetFullPath(path string) string {
	fullPath := path
	if !strings.HasPrefix(path, fileSystem.Base) {
		fullPath, _ = filepath.Abs(filepath.Join(fileSystem.Base, path))
	}
	return fullPath
}

// Get receive file with given path
func (fileSystem FileSystem) Get(path string) (*os.File, error) {
	return os.Open(fileSystem.GetFullPath(path))
}

// GetStream get file as stream
func (fileSystem FileSystem) GetStream(path string) (io.ReadCloser, error) {
	return os.Open(fileSystem.GetFullPath(path))
}

// Put store a reader into given path
func (fileSystem FileSystem) Put(path string, reader io.Reader) (*oss.Object, error) {
	var (
		fullPath = fileSystem.GetFullPath(path)
		err      = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	)

	if err != nil {
		return nil, err
	}

	dst, err := os.Create(fullPath)

	if err == nil {
		if seeker, ok := reader.(io.ReadSeeker); ok {
			_, err := seeker.Seek(0, 0)
			if err != nil {
				return nil, err
			}
		}
		_, err = io.Copy(dst, reader)
	}

	return &oss.Object{Path: path, Name: filepath.Base(path), StorageInterface: fileSystem}, err
}

// Delete delete file
func (fileSystem FileSystem) Delete(path string) error {
	return os.Remove(fileSystem.GetFullPath(path))
}

// List all objects under current path
func (fileSystem FileSystem) List(path string) ([]*oss.Object, error) {
	var (
		objects  []*oss.Object
		fullPath = fileSystem.GetFullPath(path)
	)

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if path == fullPath {
			return nil
		}

		if err == nil && !info.IsDir() {
			modTime := info.ModTime()
			objects = append(objects, &oss.Object{
				Path:             strings.TrimPrefix(path, fileSystem.Base),
				Name:             info.Name(),
				LastModified:     &modTime,
				StorageInterface: fileSystem,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// GetEndpoint get endpoint, FileSystem's endpoint is /
func (fileSystem FileSystem) GetEndpoint() string {
	return "/"
}

// GetURL get public accessible URL
func (fileSystem FileSystem) GetURL(path string) (url string, err error) {
	return path, nil
}

func NewLocalFileSystemStorageProvider() oss.StorageInterface {
	return NewFileSystem(baseFolder)
}
