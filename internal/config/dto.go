package config

import "runtime/debug"

const AppName = "git-update"

const (
	buildInfoVCSRevision = "vcs.revision"
	buildInfoVCSTime     = "vcs.time"
	unknownBuildValue    = "-"
)

// Build information -ldflags .
var (
	version    = unknownBuildValue
	branch     = unknownBuildValue
	commitHash = unknownBuildValue
	timeBuild  = unknownBuildValue
)

type Version struct {
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
	Branch     string `json:"branch,omitempty"`
	CommitHash string `json:"commitHash,omitempty"`
	TimeBuild  string `json:"timeBuild,omitempty"`
}

func GetVersion() Version {
	result := Version{
		Name:       AppName,
		Version:    version,
		Branch:     branch,
		CommitHash: commitHash,
		TimeBuild:  timeBuild,
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return result
	}

	return enrichVersionFromBuildInfo(result, buildInfo)
}

func enrichVersionFromBuildInfo(result Version, buildInfo *debug.BuildInfo) Version {
	if result.Version == unknownBuildValue && buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		result.Version = buildInfo.Main.Version
	}

	for _, setting := range buildInfo.Settings {
		if setting.Value == "" {
			continue
		}

		switch setting.Key {
		case buildInfoVCSRevision:
			if result.CommitHash == unknownBuildValue {
				result.CommitHash = setting.Value
			}
		case buildInfoVCSTime:
			if result.TimeBuild == unknownBuildValue {
				result.TimeBuild = setting.Value
			}
		}
	}

	return result
}
