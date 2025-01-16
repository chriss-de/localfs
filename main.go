package localFS

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type LocalFS struct {
	rootPath string
}

func NewLocalFS(rootPath string) (_ *LocalFS, err error) {
	var rootPathAbs string
	if rootPathAbs, err = filepath.Abs(rootPath); err != nil {
		return nil, err
	}
	return &LocalFS{rootPath: rootPathAbs}, nil
}

func (lfs *LocalFS) Open(name string) (_ fs.File, err error) {
	name = filepath.Join(lfs.rootPath, name)
	if name, err = filepath.Abs(name); err != nil {
		return nil, err
	}
	if strings.HasPrefix(name, lfs.rootPath) {
		return os.Open(name)
	}
	return nil, os.ErrNotExist
}

//func (lfs *LocalFS) ReadFile(name string) ([]byte, error) {
//	slog.Info("localFS/ReadFile", "name", name)
//	return nil, nil
//}
