package file

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"

	git "github.com/go-git/go-git/v5"
)

func testPath() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	wd := filepath.Dir(b)
	testPath := filepath.Join(wd, "..", "..", "tmp", "test")
	return testPath, nil
}

func testGitDirs() []string {
	testPath, _ := testPath()
	return []string{
		filepath.Join(testPath, "plaingitdir"),
		filepath.Join(testPath, "nocommitpushdir"),
	}
}

func testSetupGitDir() error {
	testPath, err := testPath()
	if err != nil {
		return err
	}
	gitDirs := testGitDirs()
	if err := os.RemoveAll(testPath); err != nil {
		return err
	}
	if err := os.MkdirAll(testPath, 0777); err != nil {
		return err
	}

	path := filepath.Join(testPath, "nongit")
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}
	for _, gitDir := range gitDirs {
		path := filepath.Join(testPath, gitDir)
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
		if _, err := git.PlainInit(path, false); err != nil {
			return err
		}
	}

	return nil
}
func TestFindGitDir(t *testing.T) {
	testSetupGitDir()
	testPath, err := testPath()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	f := NewFile(testPath)
	actual, err := f.FindGitDir()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	expects := testGitDirs()
	sort.Slice(actual, func(i, j int) bool { return actual[i] < actual[j] })
	sort.Slice(expects, func(i, j int) bool { return expects[i] < expects[j] })
	if !reflect.DeepEqual(actual, expects) {
		t.Errorf("expects: %v, actual: %v", expects, actual)
	}
}
