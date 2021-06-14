/*
Copyright © 2021 kouki@kmdkuk

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var finder string

// cdCmd represents the cd command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Use peco to move to the dirty Directory.",
	Long:  `Use peco to move to the dirty Directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := findDirtyGit(filePath)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}

		result, err := runFinder(dir)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		err = os.Chdir(result)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		shell, err := getShell(os.Environ())
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
		s := exec.Command(shell)
		s.Stdin = os.Stdin
		s.Stdout = os.Stdout
		s.Stderr = os.Stderr
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println("[WARNING][git-push-notifier]Exit with `exit`, since the child process is responsible for the directory change.")
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		if err := s.Run(); err != nil {
			return xerrors.Errorf("%w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)

	cdCmd.Flags().StringVarP(&finder, "finder", "f", "peco", "Select the finder you want to use()")
}

func runFinder(dirs []string) (string, error) {
	f := exec.Command(finder)
	f.Stderr = os.Stderr
	in, err := f.StdinPipe()
	if err != nil {
		return "", err
	}
	errCh := make(chan error, 1)
	go func() {
		SliceRead(dirs, in)
		errCh <- in.Close()
	}()

	err = <-errCh
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}
	result, err := f.Output()
	if err != nil {
		return "", xerrors.Errorf("%w", err)
	}

	return strings.TrimRight(string(result), "\n"), nil
}
func SliceRead(s []string, out io.WriteCloser) {
	for _, item := range s {
		fmt.Fprintln(out, item)
	}
}

func getShell(env []string) (string, error) {
	for _, s := range env {
		kv := strings.Split(s, "=")
		if kv[0] == "SHELL" {
			return kv[1], nil
		}
	}
	return "", xerrors.New("not found $SHELL")
}
