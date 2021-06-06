package file

import (
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
	err := filepath.Walk(f.root, func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, ".git")); !os.IsNotExist(err) && info.IsDir() {
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
		return nil, errors.WithStack(err)
	}
	return paths, nil
}
