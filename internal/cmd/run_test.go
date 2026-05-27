package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Format-C-eft/git-update/internal/config"
)

const testDefaultBranch = "master"

func TestRunValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr string
	}{
		{
			name: "empty path",
			setup: func(t *testing.T) {
				restoreConfig(t)

				config.FlagDir = ""
				config.FlagCheckout = true
			},
			wantErr: "empty path to directory",
		},
		{
			name: "actions disabled",
			setup: func(t *testing.T) {
				restoreConfig(t)

				config.FlagDir = t.TempDir()
			},
			wantErr: "all actions are disabled",
		},
		{
			name: "invalid parallel",
			setup: func(t *testing.T) {
				restoreConfig(t)

				config.FlagDir = t.TempDir()
				config.FlagCheckout = true
				config.FlagParallel = 0
			},
			wantErr: "parallel must be greater than zero",
		},
		{
			name: "no subdirectories",
			setup: func(t *testing.T) {
				restoreConfig(t)

				config.FlagDir = t.TempDir()
				config.FlagCheckout = true
			},
			wantErr: "directory does not contain subdirectories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			err := Run()
			if err == nil {
				t.Fatalf("Run returned nil error, want %q", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Run error = %q, want to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestGetListOfDirectories(t *testing.T) {
	restoreConfig(t)

	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "repo-a"))
	mustMkdir(t, filepath.Join(root, "repo-b"))
	mustWriteFile(t, filepath.Join(root, "README.md"), "not a directory")

	config.FlagDir = root + string(os.PathSeparator)

	got, err := getListOfDirectories()
	if err != nil {
		t.Fatalf("getListOfDirectories returned error: %v", err)
	}

	sort.Strings(got)

	want := []string{
		filepath.Join(root, "repo-a"),
		filepath.Join(root, "repo-b"),
	}

	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("directories = %#v, want %#v", got, want)
	}
}

func TestRunCheckoutCleanRepository(t *testing.T) {
	restoreConfig(t)
	requireGit(t)

	root := t.TempDir()
	repo := filepath.Join(root, testRepoName)
	mustMkdir(t, repo)
	initGitRepo(t, repo)

	config.FlagDir = root
	config.FlagCheckout = true
	config.FlagDefaultBranch = testDefaultBranch

	err := captureStdout(t, Run)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestProcessDirSkipsDirtyRepositoryWithoutReset(t *testing.T) {
	restoreConfig(t)
	requireGit(t)

	repo := filepath.Join(t.TempDir(), testRepoName)
	mustMkdir(t, repo)
	initGitRepo(t, repo)
	mustWriteFile(t, filepath.Join(repo, "dirty.txt"), "dirty")

	config.FlagCheckout = true
	config.FlagDefaultBranch = testDefaultBranch
	config.FlagResetHard = false

	result := runProcessDir(t, repo)
	got := result.String()

	if !strings.Contains(got, "find uncommitted changes") {
		t.Fatalf("processDir log = %q, want dirty repository message", got)
	}

	if !strings.Contains(got, "skipped: there are uncommitted changes") {
		t.Fatalf("processDir log = %q, want skip message", got)
	}
}

func TestProcessDirPullFastForwardOnlyUpdatesRepository(t *testing.T) {
	restoreConfig(t)

	_, seed, repo := createRemoteWithClone(t)

	writeFileAndCommit(t, seed, "file.txt", "remote update\n", "remote update")
	runGit(t, seed, "push", "origin", testDefaultBranch)

	config.FlagPull = true

	result := runProcessDir(t, repo)

	got := result.String()
	if !strings.Contains(got, "success: git pull --ff-only") {
		t.Fatalf("processDir log = %q, want successful fast-forward pull", got)
	}

	//nolint:gosec // Test reads a file from a controlled temporary repository.
	content, err := os.ReadFile(filepath.Join(repo, "file.txt"))
	if err != nil {
		t.Fatalf("read pulled file: %v", err)
	}

	if string(content) != "remote update\n" {
		t.Fatalf("pulled file content = %q, want %q", string(content), "remote update\n")
	}
}

func TestProcessDirPullFastForwardOnlyRejectsDivergedHistory(t *testing.T) {
	restoreConfig(t)

	_, seed, repo := createRemoteWithClone(t)

	writeFileAndCommit(t, seed, "remote.txt", "remote\n", "remote update")
	runGit(t, seed, "push", "origin", testDefaultBranch)
	writeFileAndCommit(t, repo, "local.txt", "local\n", "local update")

	config.FlagPull = true

	result := runProcessDir(t, repo)

	got := result.String()
	if !strings.Contains(got, "error: git pull --ff-only") {
		t.Fatalf("processDir log = %q, want fast-forward-only error", got)
	}
}

func restoreConfig(t *testing.T) {
	t.Helper()

	oldDir := config.FlagDir
	oldAll := config.FlagAll
	oldCheckout := config.FlagCheckout
	oldFetch := config.FlagFetch
	oldPull := config.FlagPull
	oldResetHard := config.FlagResetHard
	oldDefaultBranch := config.FlagDefaultBranch
	oldVerbose := config.FlagVerbose
	oldParallel := config.FlagParallel
	oldExecuteTimeout := config.FlagExecuteTimeout
	oldVersion := config.FlagVersion

	t.Cleanup(func() {
		config.FlagDir = oldDir
		config.FlagAll = oldAll
		config.FlagCheckout = oldCheckout
		config.FlagFetch = oldFetch
		config.FlagPull = oldPull
		config.FlagResetHard = oldResetHard
		config.FlagDefaultBranch = oldDefaultBranch
		config.FlagVerbose = oldVerbose
		config.FlagParallel = oldParallel
		config.FlagExecuteTimeout = oldExecuteTimeout
		config.FlagVersion = oldVersion
	})

	config.FlagDir = ""
	config.FlagAll = false
	config.FlagCheckout = false
	config.FlagFetch = false
	config.FlagPull = false
	config.FlagResetHard = false
	config.FlagDefaultBranch = testDefaultBranch
	config.FlagVerbose = false
	config.FlagParallel = 1
	config.FlagExecuteTimeout = 5 * time.Second
	config.FlagVersion = false
}

func captureStdout(t *testing.T, fn func() error) error {
	t.Helper()

	oldStdout := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = oldStdout
	}()

	errRun := fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}

	if _, err := io.ReadAll(reader); err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}

	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}

	return errRun
}

func runProcessDir(t *testing.T, repo string) resultLog {
	t.Helper()

	ch := make(chan resultLog, 1)
	processDir(repo, ch)

	return <-ch
}

func createRemoteWithClone(t *testing.T) (string, string, string) {
	t.Helper()
	requireGit(t)

	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	repo := filepath.Join(root, testRepoName)

	mustMkdir(t, remote)
	mustMkdir(t, seed)

	runGit(t, remote, "init", "--bare")
	initGitRepo(t, seed)
	writeFileAndCommit(t, seed, "file.txt", "base\n", "initial commit")
	runGit(t, seed, "remote", "add", "origin", remote)
	runGit(t, seed, "push", "-u", "origin", testDefaultBranch)
	runGit(t, remote, "symbolic-ref", "HEAD", "refs/heads/"+testDefaultBranch)
	runGit(t, root, "clone", remote, repo)

	return remote, seed, repo
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "init", "-b", testDefaultBranch)
}

func writeFileAndCommit(t *testing.T, dir string, name string, content string, message string) {
	t.Helper()

	mustWriteFile(t, filepath.Join(dir, name), content)
	runGit(t, dir, "add", name)
	runGit(t, dir,
		"-c", "user.name=Test User",
		"-c", "user.email=test@example.com",
		"commit",
		"-m", message,
	)
}

func requireGit(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is required for this test")
	}
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	gitArgs := append([]string{"-C", dir}, args...)
	//nolint:gosec // Tests execute git with controlled arguments.
	cmd := exec.Command("git", gitArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(gitArgs, " "), err, output)
	}

	return string(output)
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o750); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
