package config

const AppName = "git-update"

// Build information -ldflags .
var (
	branch     = "dev"
	commitHash = "-"
	timeBuild  = "-"
)

type Version struct {
	Name       string `json:"name,omitempty"`
	Branch     string `json:"branch,omitempty"`
	CommitHash string `json:"commitHash,omitempty"`
	TimeBuild  string `json:"timeBuild,omitempty"`
}

func GetVersion() Version {
	return Version{
		Name:       AppName,
		Branch:     branch,
		CommitHash: commitHash,
		TimeBuild:  timeBuild,
	}
}
