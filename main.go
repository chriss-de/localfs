package localFS

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type LocalFS struct {
	rootPath string
	tryFiles []string
}

func NewLocalFS(rootPath string, opts ...func(lfs *LocalFS)) (lfs *LocalFS, err error) {
	var rootPathAbs string
	if rootPathAbs, err = filepath.Abs(rootPath); err != nil {
		return nil, err
	}

	lfs = &LocalFS{rootPath: rootPathAbs}

	for _, opt := range opts {
		opt(lfs)
	}

	return lfs, nil
}

func WithTryFile(file string) func(lfs *LocalFS) {
	return func(lfs *LocalFS) {
		lfs.tryFiles = append(lfs.tryFiles, file)
	}
}

func WithTryFiles(files ...string) func(lfs *LocalFS) {
	return func(lfs *LocalFS) {
		lfs.tryFiles = append(lfs.tryFiles, files...)
	}
}

func (lfs *LocalFS) Open(name string) (fd fs.File, err error) {

	var listOfFiles = []string{"name"}
	listOfFiles = append(listOfFiles, lfs.tryFiles...)

	return lfs.tryOpen(listOfFiles...)
}

func (lfs *LocalFS) tryOpen(fileNames ...string) (fd fs.File, err error) {
	for _, filename := range fileNames {
		filename = filepath.Join(lfs.rootPath, filename)
		if filename, err = filepath.Abs(filename); err != nil {
			return nil, err
		}
		if strings.HasPrefix(filename, lfs.rootPath) {
			if fd, err = os.Open(filename); err == nil {
				return fd, nil
			}
		}
	}

	return nil, os.ErrNotExist
}
