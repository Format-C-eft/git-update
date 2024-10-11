package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Format-C-eft/git-update/internal/config"
	"github.com/Format-C-eft/git-update/internal/shell_executor"
)

func Run() error {
	if config.FlagDir == "" {
		return errors.New("empty path to directory")
	}

	if !strings.HasSuffix(config.FlagDir, string(os.PathSeparator)) {
		config.FlagDir += string(os.PathSeparator)
	}

	if config.FlagAll {
		config.FlagCheckout = true
		config.FlagFetch = true
		config.FlagPull = true
	}

	if !config.FlagCheckout && !config.FlagFetch && !config.FlagPull {
		return errors.New("all actions are disabled")
	}

	listDir, err := getListOfDirectories()
	if err != nil {
		return fmt.Errorf("getListOfDirectories err: %w", err)
	}

	if len(listDir) == 0 {
		return errors.New("directory does not contain subdirectories")
	}

	ch := make(chan resultLog, len(listDir))
	defer close(ch)

	for _, dir := range listDir {
		go processDir(dir, ch)
	}

	execTimeout := *config.FlagExecuteTimeout
	execTimeout = time.Duration(int64(len(listDir)) * execTimeout.Nanoseconds())

	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	var countDone atomic.Int32
	var countAll = int32(len(listDir))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case logChan := <-ch:
			fmt.Println(logChan.String())
			fmt.Println("--------------------------------------------------------------------")
			if newValue := countDone.Add(1); newValue >= countAll {
				return nil
			}
		}
	}
}

func getListOfDirectories() ([]string, error) {
	list, err := os.ReadDir(config.FlagDir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir err: %w", err)
	}

	listDir := make([]string, 0, len(list))
	for _, path := range list {
		if path.IsDir() {
			listDir = append(listDir, fmt.Sprintf("%s%s", config.FlagDir, path.Name()))
		}
	}

	return listDir, nil
}

func processDir(dir string, ch chan resultLog) {
	resultLogs := resultLog{dir: dir}
	resultLogs.AddLog("start processing", "")

	defer func() {
		resultLogs.AddLog("stop processing", "")
		ch <- resultLogs
	}()

	ctx, cancelFn := context.WithTimeout(context.Background(), *config.FlagExecuteTimeout)
	defer cancelFn()

	res, errRun := shell_executor.Run(ctx, dir, "git", "status")
	if errRun != nil {
		resultLogs.AddLog("skipped: error: execute git status", errRun.Error())
		return
	}

	if !strings.Contains(res, gitStatusOk) {
		resultLogs.AddLog("find uncommitted changes", "")

		if !config.FlagResetHard {
			resultLogs.AddLog("skipped: there are uncommitted changes", "")
			return
		}
		result, err := shell_executor.Run(ctx, dir, "git", "reset", "--hard")
		if err != nil {
			resultLogs.AddLog("error: git reset --hard", err.Error())
			return
		}

		resultLogs.AddLog("success: git reset --hard", result)
	}

	if config.FlagCheckout {
		result, err := shell_executor.Run(ctx, dir, "git", "checkout", config.FlagDefaultBranch)
		if err != nil {
			resultLogs.AddLog("error: git checkout", err.Error())
			return
		}

		resultLogs.AddLog("success: git checkout", result)
	}

	if config.FlagFetch {
		result, err := shell_executor.Run(ctx, dir, "git", "fetch", "--prune", "--prune-tags")
		if err != nil {
			resultLogs.AddLog("error: git fetch", err.Error())
			return
		}

		resultLogs.AddLog("success: git fetch", result)
	}

	if config.FlagPull {
		result, err := shell_executor.Run(ctx, dir, "git", "pull")
		if err != nil {
			resultLogs.AddLog("error: git pull ", err.Error())
			return
		}

		resultLogs.AddLog("success: git pull ", result)
	}
}
