package git

type Git struct {
	// 対象のgitrepositoryの絶対パスリスト
	paths []string
}

func NewGit(paths []string) *Git {
	return &Git{
		paths: paths,
	}
}

func (g *Git) FindDirtyGit() ([]string, error) {
	return g.paths, nil
}
