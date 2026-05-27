package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainWritesRunErrorToStderrAndExitsNonZero(t *testing.T) {
	if os.Getenv("GO_WANT_MAIN_HELPER_PROCESS") == "1" {
		main()

		return
	}

	//nolint:gosec // Test executes the current test binary as a helper process.
	cmd := exec.Command(os.Args[0], "-test.run=TestMainWritesRunErrorToStderrAndExitsNonZero")

	cmd.Env = append(os.Environ(), "GO_WANT_MAIN_HELPER_PROCESS=1")

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Fatal("helper process exited successfully, want failure")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("helper process error = %v, want exec.ExitError", err)
	}

	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want %d", exitErr.ExitCode(), 1)
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	if !strings.Contains(stderr.String(), "all actions are disabled") {
		t.Fatalf("stderr = %q, want run error", stderr.String())
	}
}
