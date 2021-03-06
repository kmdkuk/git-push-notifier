package file

import (
	"reflect"
	"testing"

	. "github.com/kmdkuk/git-push-notifier/helper/test"
)

func TestFindGitDir(t *testing.T) {
	testPath := GetPath()
	f := NewFile(testPath)
	actual, err := f.FindGitDir()
	if err != nil {
		t.Errorf("error: %+v", err)
	}
	expects := GetGitDirs()
	SortString(actual)
	SortString(expects)
	if !reflect.DeepEqual(actual, expects) {
		t.Errorf("expects: %+v, actual: %+v", expects, actual)
	}
}
