package file

import "fmt"

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
	return nil, fmt.Errorf("not find git dir")
}
