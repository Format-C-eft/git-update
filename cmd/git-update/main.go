package main

import (
	"fmt"
	"os"

	"github.com/Format-C-eft/git-update/internal/cmd"
	"github.com/Format-C-eft/git-update/internal/config"
)

func main() {
	config.ParseFlags()

	if config.FlagVersion {
		showVersion()
		return
	}

	if errRun := cmd.Run(); errRun != nil {
		_, _ = fmt.Fprintln(os.Stderr, errRun.Error())

		os.Exit(1)
	}
}

func showVersion() {
	version := config.GetVersion()

	fmt.Printf("Name - '%s'\n", version.Name)
	fmt.Printf("Version - '%s'\n", version.Version)
	fmt.Printf("Branch - '%s'\n", version.Branch)
	fmt.Printf("Commit hash - '%s'\n", version.CommitHash)
	fmt.Printf("Time build - '%s'\n", version.TimeBuild)
}
