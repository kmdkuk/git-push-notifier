package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/kmdkuk/git-push-notifier/log"
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

			if isNotCommit(dir) || isNotPush(dir) {
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

func isNotCommit(dir string) bool {
	status, err := runCommand("git", "-C", dir, "status", "--porcelain")
	if err != nil {
		return false
	}
	return len(status) != 0
}

func isNotPush(dir string) bool {
	branchs, err := getBranchList(dir)
	if err != nil {
		return false
	}
	for _, branch := range branchs {
		if !isExistOrigin(dir, branch) {
			log.Debugf("is not found origin/%s in %s", branch, dir)
			return true
		}

		logdiff, err := runCommand("git", "-C", dir, "log", fmt.Sprintf("origin/%s..%s", branch, branch))
		if err != nil {
			continue
		}
		if len(logdiff) != 0 {
			log.Debugf("not push %s to origin/%s in %s", branch, branch, dir)
			return true
		}
	}
	return false
}

func isExistOrigin(dir string, branch string) bool {
	remoteBranchsStr, err := runCommand("git", "-C", dir, "branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return false
	}
	remoteBranchs := strings.Split(remoteBranchsStr, "\n")
	for _, remoteBranch := range remoteBranchs {
		if strings.TrimPrefix(remoteBranch, "origin/") == branch {
			return true
		}
	}
	return false
}

func getBranchList(dir string) ([]string, error) {
	branchsStr, err := runCommand("git", "-C", dir, "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	branchs := strings.Split(branchsStr, "\n")
	result := make([]string, 0)
	for _, b := range branchs {
		log.Debug(b)
		if strings.Contains(b, "HEAD detached at") {
			continue
		}
		trimed := strings.TrimSpace(b)
		if len(trimed) != 0 {
			result = append(result, trimed)
		}
	}

	return result, nil
}

func runCommand(cmdName string, args ...string) (string, error) {
	var buffer bytes.Buffer
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = &buffer
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
