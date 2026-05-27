package config

import (
	"runtime/debug"
	"testing"
	"time"
)

const (
	testCommitHash      = "abc123"
	testModuleVersion   = "v0.1.1"
	testOverrideVersion = "v9.9.9"
	testRevision        = "1234567890abcdef"
	testTimeBuild       = "2026-05-27T12:00:00+0500"
	testVCSTime         = "2026-05-27T08:00:00Z"
)

func TestGetVersion(t *testing.T) {
	oldVersion := version
	oldBranch := branch
	oldCommitHash := commitHash
	oldTimeBuild := timeBuild

	t.Cleanup(func() {
		version = oldVersion
		branch = oldBranch
		commitHash = oldCommitHash
		timeBuild = oldTimeBuild
	})

	version = "v1.2.3"
	branch = "test-branch"
	commitHash = testCommitHash
	timeBuild = testTimeBuild

	got := GetVersion()

	if got.Name != AppName {
		t.Fatalf("Name = %q, want %q", got.Name, AppName)
	}

	if got.Version != version {
		t.Fatalf("Version = %q, want %q", got.Version, version)
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

func TestEnrichVersionFromBuildInfo(t *testing.T) {
	got := enrichVersionFromBuildInfo(Version{
		Name:       AppName,
		Version:    unknownBuildValue,
		Branch:     unknownBuildValue,
		CommitHash: unknownBuildValue,
		TimeBuild:  unknownBuildValue,
	}, &debug.BuildInfo{
		Main: debug.Module{
			Version: testModuleVersion,
		},
		Settings: []debug.BuildSetting{
			{Key: buildInfoVCSRevision, Value: testRevision},
			{Key: buildInfoVCSTime, Value: testVCSTime},
		},
	})

	if got.Version != testModuleVersion {
		t.Fatalf("Version = %q, want %q", got.Version, testModuleVersion)
	}

	if got.CommitHash != testRevision {
		t.Fatalf("CommitHash = %q, want %q", got.CommitHash, testRevision)
	}

	if got.TimeBuild != testVCSTime {
		t.Fatalf("TimeBuild = %q, want %q", got.TimeBuild, testVCSTime)
	}
}

func TestEnrichVersionFromBuildInfoKeepsLdflagsValues(t *testing.T) {
	got := enrichVersionFromBuildInfo(Version{
		Name:       AppName,
		Version:    testOverrideVersion,
		Branch:     "main",
		CommitHash: testCommitHash,
		TimeBuild:  testTimeBuild,
	}, &debug.BuildInfo{
		Main: debug.Module{
			Version: testModuleVersion,
		},
		Settings: []debug.BuildSetting{
			{Key: buildInfoVCSRevision, Value: testRevision},
			{Key: buildInfoVCSTime, Value: testVCSTime},
		},
	})

	if got.Version != testOverrideVersion {
		t.Fatalf("Version = %q, want %q", got.Version, testOverrideVersion)
	}

	if got.CommitHash != testCommitHash {
		t.Fatalf("CommitHash = %q, want %q", got.CommitHash, testCommitHash)
	}

	if got.TimeBuild != testTimeBuild {
		t.Fatalf("TimeBuild = %q, want %q", got.TimeBuild, testTimeBuild)
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
