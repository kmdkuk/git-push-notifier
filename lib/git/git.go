package git

import (
	"sync"

	gitv5 "github.com/go-git/go-git/v5"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

type Git struct {
	sync.Mutex
	// 対象のgitrepositoryの絶対パスリスト
	paths []string
}

func NewGit(paths []string) *Git {
	return &Git{
		paths: paths,
	}
}

func (g *Git) FindDirtyGit() ([]string, error) {
	dirtyDir := make([]string, 0)
	eg := errgroup.Group{}
	for _, dir := range g.paths {
		dir := dir
		eg.Go(func() error {
			r, err := gitv5.PlainOpen(dir)
			if err != nil {
				return xerrors.Errorf("%w", err)
			}

			w, err := r.Worktree()
			if err != nil {
				return xerrors.Errorf("%w", err)
			}

			status, err := w.Status()
			if err != nil {
				return xerrors.Errorf("%w", err)
			}

			g.Lock()
			if !status.IsClean() {
				dirtyDir = append(dirtyDir, dir)
			}
			g.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return dirtyDir, nil
}
