package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/Format-C-eft/git-update/internal/config"
)

const testRepoName = "repo"

func TestResultLogAddLogRespectsVerboseFlag(t *testing.T) {
	restoreConfig(t)

	result := resultLog{dir: testRepoName}

	config.FlagVerbose = false

	result.AddLog("plain", "hidden")

	if len(result.logs) != 1 {
		t.Fatalf("logs length = %d, want %d", len(result.logs), 1)
	}

	if result.logs[0].messages != "plain" {
		t.Fatalf("first log message = %q, want %q", result.logs[0].messages, "plain")
	}

	config.FlagVerbose = true

	result.AddLog("plain again", "visible")

	if len(result.logs) != 3 {
		t.Fatalf("logs length = %d, want %d", len(result.logs), 3)
	}

	if result.logs[2].messages != "verbose: visible" {
		t.Fatalf("verbose log message = %q, want %q", result.logs[2].messages, "verbose: visible")
	}
}

func TestResultLogString(t *testing.T) {
	first := time.Date(2026, 5, 27, 12, 34, 56, 789*int(time.Millisecond), time.UTC)
	second := first.Add(time.Millisecond)
	result := resultLog{
		dir: testRepoName,
		logs: []logs{
			{Time: first, messages: "start"},
			{Time: second, messages: "stop"},
		},
	}

	got := result.String()
	want := strings.Join([]string{
		"12:34:56.789: " + testRepoName + ": start",
		"12:34:56.79: " + testRepoName + ": stop",
	}, "\n")

	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
