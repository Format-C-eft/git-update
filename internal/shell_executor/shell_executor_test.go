package shell_executor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunReturnsSanitizedStdout(t *testing.T) {
	got, err := Run(context.Background(), t.TempDir(), os.Args[0], "-test.run=TestHelperProcess", "--", "stdout")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	want := "line one line two "
	if got != want {
		t.Fatalf("Run stdout = %q, want %q", got, want)
	}
}

func TestRunReturnsSanitizedStderrOnFailure(t *testing.T) {
	_, err := Run(context.Background(), t.TempDir(), os.Args[0], "-test.run=TestHelperProcess", "--", "stderr")
	if err == nil {
		t.Fatal("Run returned nil error, want failure")
	}

	if !strings.Contains(err.Error(), "bad news ") {
		t.Fatalf("Run error = %q, want sanitized stderr", err.Error())
	}
}

func TestRunHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := Run(ctx, t.TempDir(), os.Args[0], "-test.run=TestHelperProcess", "--", "sleep")
	if err == nil {
		t.Fatal("Run returned nil error, want context cancellation failure")
	}
}

func TestRunReportsStartError(t *testing.T) {
	missingBinary := t.TempDir() + string(os.PathSeparator) + "missing-binary"

	_, err := Run(context.Background(), t.TempDir(), missingBinary)
	if err == nil {
		t.Fatal("Run returned nil error, want start failure")
	}

	if !strings.Contains(err.Error(), missingBinary) {
		t.Fatalf("Run error = %q, want missing binary path", err.Error())
	}
}

func TestHelperProcess(t *testing.T) {
	args := os.Args
	for i, arg := range args {
		if arg != "--" {
			continue
		}

		runHelperProcess(args[i+1:])
	}
}

func runHelperProcess(args []string) {
	if len(args) == 0 {
		os.Exit(0)
	}

	switch args[0] {
	case "stdout":
		fmt.Print("line one\nline two\n")
	case "stderr":
		_, _ = fmt.Fprint(os.Stderr, "bad\nnews\n")

		os.Exit(42)
	case "sleep":
		time.Sleep(time.Second)
	}

	os.Exit(0)
}
