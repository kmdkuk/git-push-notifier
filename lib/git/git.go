package git

import (
	"bytes"
	"os"
	"os/exec"
	"sync"

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
			/*
				// TODO: too slow
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
			*/

			var buffer bytes.Buffer
			statusCmd := exec.Command("git", "-C", dir, "status", "--porcelain")
			statusCmd.Stdout = &buffer
			statusCmd.Stderr = os.Stderr
			if err := statusCmd.Run(); err != nil {
				return xerrors.Errorf("%w", err)
			}
			status := buffer.String()

			if len(status) != 0 {
				g.Lock()
				dirtyDir = append(dirtyDir, dir)
				g.Unlock()
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return dirtyDir, nil
}
