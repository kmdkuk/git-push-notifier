package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type File struct {
	// pushしてないリポジトリを探すルートpath
	root string
}

func NewFile(root string) *File {
	return &File{
		root: root,
	}
}

func (f *File) FindGitDir() ([]string, error) {
	paths := make([]string, 0)
	err := filepath.WalkDir(f.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if _, err := os.Stat(filepath.Join(path, ".git")); !os.IsNotExist(err) {
			apath, err := filepath.Abs(path)
			if err != nil {
				return errors.Wrapf(err, "path: %s", path)
			}
			paths = append(paths, apath)
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return paths, nil
}
