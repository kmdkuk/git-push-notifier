package git

import (
	"reflect"
	"testing"

	. "github.com/kmdkuk/git-push-notifier/helper/test"
)

func TestFindNonStaginGit(t *testing.T) {
	gitDirs := GetGitDirs()
	g := NewGit(gitDirs)
	actual, err := g.FindDirtyGit()
	if err != nil {
		t.Errorf("error: %+v\n", err)
	}
	expects := GetNonStagingGitDirs()

	SortString(actual)
	SortString(expects)
	if !reflect.DeepEqual(actual, expects) {
		t.Errorf("expects: %+v, actual: %+v", expects, actual)
	}
}
