package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kmdkuk/git-push-notifier/log"
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
			return err
		}
		if _, err := os.Stat(filepath.Join(path, ".git")); !os.IsNotExist(err) {
			log.Debugf("root %s\n", f.root)
			log.Debugf("path %s\n", path)
			rel, err := filepath.Rel(f.root, path)
			if err != nil {
				return err
			}
			paths = append(paths, filepath.Join("/", rel))
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return paths, nil
}
