package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CmdVersion struct {
	AppVersion string
	BuildTime  string
	GitCommit  string
	GitRef     string
}

var (
	shortVersion bool
	Version      *CmdVersion

	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Displays version information",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			if shortVersion {
				fmt.Printf("%s\n", Version.AppVersion)
			} else {
				fmt.Printf("Docker Control Version:: %s\nBuildtime: %s\nGitCommit: %s\n", Version.AppVersion, Version.BuildTime, Version.GitCommit)
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&shortVersion, "short", "s", false, "Show the version string only")
}
