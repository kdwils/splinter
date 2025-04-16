package fio

import (
	"io"
	"io/fs"
	"os"
)

// FileIO is an interface for os stdlib file operations
type FileIO interface {
	ReadFile(filename string) ([]byte, error)
	Stat(name string) (fs.FileInfo, error)
	ReadDir(name string) ([]fs.DirEntry, error)
	Create(name string) (io.WriteCloser, error)
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(filename string, data []byte, perm fs.FileMode) error
}

// WriteCloser is an alias for io.WriteCloser
type WriteCloser io.WriteCloser

// DirEntry is an alias for fs.DirEntry
type DirEntry fs.DirEntry

// DefaultFS is a struct that implements the FileIO interface using the os package
type DefaultFS struct{}

// NewDefaultFileIO creates a new instance of DefaultFS that implements the FileIO interface
func NewDefaultFileIO() FileIO {
	return DefaultFS{}
}

// ReadFile is a wrapper around os.ReadFile
func (fsys DefaultFS) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// Stat is a wrapper around os.Stat
func (fsys DefaultFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// ReadDir is a wrapper around os.ReadDir
func (fsys DefaultFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

// WriteFile is a wrapper around os.WriteFile
func (fsys DefaultFS) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// Create is a wrapper around os.Create
func (fsys DefaultFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

// MkdirAll is a wrapper around os.MkdirAll
func (fsys DefaultFS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}
