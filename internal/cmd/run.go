package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Format-C-eft/git-update/internal/shell_executor"
)

func cmdRun(_ *cobra.Command, _ []string) error {
	if flagDir == "" {
		return errors.New("empty path to directory")
	}

	if !strings.HasSuffix(flagDir, string(os.PathSeparator)) {
		flagDir = flagDir + string(os.PathSeparator)
	}

	if flagAll {
		flagCheckout = true
		flagFetch = true
		flagPull = true
	}

	if !flagCheckout && !flagFetch && !flagPull {
		return errors.New("all actions are disabled")
	}

	listDir, err := getListOfDirectories()
	if err != nil {
		return errors.Wrap(err, "getListOfDirectories")
	}

	if len(listDir) == 0 {
		return errors.New("directory does not contain subdirectories")
	}

	ch := make(chan resultLogChan, len(listDir))

	go runProcess(listDir, ch)

	for logChan := range ch {
		fmt.Println(logChan.String())
	}

	return nil
}

func getListOfDirectories() ([]string, error) {
	list, err := os.ReadDir(flagDir)
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadDir")
	}

	listDir := make([]string, 0, len(list))
	for _, path := range list {
		if path.IsDir() {

			listDir = append(listDir, fmt.Sprintf("%s%s", flagDir, path.Name()))
		}
	}

	return listDir, nil
}

func runProcess(listDir []string, ch chan resultLogChan) {
	defer close(ch)

	wg := &sync.WaitGroup{}
	for _, dir := range listDir {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			processDir(dir, ch)
		}(dir)
	}

	wg.Wait()
}

func processDir(dir string, ch chan resultLogChan) {
	resultLogs := resultLogChan{dir: dir}
	resultLogs.AddLog("start processing")

	defer func() {
		resultLogs.AddLog("stop processing")
		ch <- resultLogs
	}()

	ctx, cancelFn := context.WithTimeout(context.Background(), *flagExecuteTimeout)
	defer cancelFn()

	res, errRun := shell_executor.Run(ctx, dir, "git", "status")
	if errRun != nil {
		resultLogs.AddLog("skipped: error: execute git status")
		return
	}

	if !strings.Contains(string(res), gitStatusOk) {
		resultLogs.AddLog("find uncommitted changes")

		if !flagResetHard {
			resultLogs.AddLog("skipped: there are uncommitted changes")
			return
		}

		if _, err := shell_executor.Run(ctx, dir, "git", "reset", "--hard"); err != nil {
			resultLogs.AddLog("error: git reset --hard: " + err.Error())
			return
		}
		resultLogs.AddLog("success: git reset --hard")
	}

	if _, err := shell_executor.Run(ctx, dir, "git", "checkout", flagDefaultBranch); err != nil {
		resultLogs.AddLog("error: git checkout: " + err.Error())
		return
	}
	resultLogs.AddLog("success: git checkout")

	if _, err := shell_executor.Run(ctx, dir, "git", "fetch", "--prune"); err != nil {
		resultLogs.AddLog("error: git fetch: " + err.Error())
		return
	}
	resultLogs.AddLog("success: git fetch")

	if _, err := shell_executor.Run(ctx, dir, "git", "pull"); err != nil {
		resultLogs.AddLog("error: git pull: " + err.Error())
		return
	}
	resultLogs.AddLog("success: git pull")

}
