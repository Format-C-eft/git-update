package config

import (
	"testing"
	"time"
)

func TestGetVersion(t *testing.T) {
	oldBranch := branch
	oldCommitHash := commitHash
	oldTimeBuild := timeBuild

	t.Cleanup(func() {
		branch = oldBranch
		commitHash = oldCommitHash
		timeBuild = oldTimeBuild
	})

	branch = "test-branch"
	commitHash = "abc123"
	timeBuild = "2026-05-27T12:00:00+0500"

	got := GetVersion()

	if got.Name != AppName {
		t.Fatalf("Name = %q, want %q", got.Name, AppName)
	}

	if got.Branch != branch {
		t.Fatalf("Branch = %q, want %q", got.Branch, branch)
	}

	if got.CommitHash != commitHash {
		t.Fatalf("CommitHash = %q, want %q", got.CommitHash, commitHash)
	}

	if got.TimeBuild != timeBuild {
		t.Fatalf("TimeBuild = %q, want %q", got.TimeBuild, timeBuild)
	}
}

func TestDefaultFlags(t *testing.T) {
	if FlagDir != "../" {
		t.Fatalf("FlagDir = %q, want %q", FlagDir, "../")
	}

	if FlagDefaultBranch != "master" {
		t.Fatalf("FlagDefaultBranch = %q, want %q", FlagDefaultBranch, "master")
	}

	if FlagParallel != 4 {
		t.Fatalf("FlagParallel = %d, want %d", FlagParallel, 4)
	}

	if FlagExecuteTimeout != 30*time.Second {
		t.Fatalf("FlagExecuteTimeout = %s, want %s", FlagExecuteTimeout, 30*time.Second)
	}
}
