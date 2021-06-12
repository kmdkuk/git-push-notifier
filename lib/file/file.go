package file

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
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
				return xerrors.Errorf("path = %s, err: %w", path, err)
			}
			paths = append(paths, apath)
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return paths, nil
}
