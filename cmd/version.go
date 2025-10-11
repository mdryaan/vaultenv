package cmd

import (
	"fmt"

	"github.com/mdryaan/vaultenv/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version, git commit, build date, and Go runtime information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
