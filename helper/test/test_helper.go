package test

import (
	"path/filepath"
	"runtime"
	"sort"
)

func GetPath() string {
	_, b, _, _ := runtime.Caller(0)
	wd := filepath.Dir(b)
	testPath := filepath.Join(wd, "..", "..", "tmp", "test")
	return testPath
}

func GetGitDirs() []string {
	dirAll := make([]string, 0)
	dirAll = append(dirAll, GetCleanGitDirs()...)
	dirAll = append(dirAll, GetNonStagingGitDirs()...)
	return dirAll
}

func GetCleanGitDirs() []string {
	testPath := GetPath()
	return []string{
		filepath.Join(testPath, "plaingitdir"),
	}
}

func GetNonStagingGitDirs() []string {
	testPath := GetPath()
	return []string{
		filepath.Join(testPath, "nocommitpushdir"),
	}
}

func SortString(strs []string) {
	sort.Slice(strs, func(i, j int) bool { return strs[i] < strs[j] })
}
