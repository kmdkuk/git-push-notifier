package file

import (
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"
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

func TestFindGitDir(t *testing.T) {
	testPath, err := testPath()
	if err != nil {
		t.Errorf("error: %+v", err)
	}
	f := NewFile(testPath)
	actual, err := f.FindGitDir()
	if err != nil {
		t.Errorf("error: %+v", err)
	}
	expects := testGitDirs()
	sort.Slice(actual, func(i, j int) bool { return actual[i] < actual[j] })
	sort.Slice(expects, func(i, j int) bool { return expects[i] < expects[j] })
	if !reflect.DeepEqual(actual, expects) {
		t.Errorf("expects: %+v, actual: %+v", expects, actual)
	}
}
