package main

import (
	"fmt"

	"github.com/Format-C-eft/git-update/internal/cmd"
	"github.com/Format-C-eft/git-update/internal/config"
)

func main() {
	if config.FlagVersion {
		showVersion()
		return
	}

	if errRun := cmd.Run(); errRun != nil {
		fmt.Println(errRun.Error())
		return
	}
}

func showVersion() {
	fmt.Printf("Name - '%s'\n", config.GetVersion().Name)
	fmt.Printf("Branch - '%s'\n", config.GetVersion().Branch)
	fmt.Printf("Commit hash - '%s'\n", config.GetVersion().CommitHash)
	fmt.Printf("Time build - '%s'\n", config.GetVersion().TimeBuild)
}
