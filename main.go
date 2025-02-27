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

func WithTryFile(files string) func(lfs *LocalFS) {
	return func(lfs *LocalFS) {
		lfs.tryFiles = append(lfs.tryFiles, files)
	}
}

func WithTryFiles(files ...string) func(lfs *LocalFS) {
	return func(lfs *LocalFS) {
		lfs.tryFiles = append(lfs.tryFiles, files...)
	}
}

func (lfs *LocalFS) Open(name string) (fd fs.File, err error) {
	name = filepath.Join(lfs.rootPath, name)
	if name, err = filepath.Abs(name); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(name, lfs.rootPath) {
		return nil, os.ErrNotExist
	}

	return lfs.tryOpen(name)
}

func (lfs *LocalFS) tryOpen(name string) (fd fs.File, err error) {
	if fd, err = os.Open(name); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		for _, tryFile := range lfs.tryFiles {
			if fd, err = os.Open(tryFile); err != nil {
				if !os.IsNotExist(err) {
					return nil, err
				}
			} else {
				break
			}
		}

		if fd == nil {
			return nil, os.ErrNotExist
		}
	}

	return fd, nil
}
