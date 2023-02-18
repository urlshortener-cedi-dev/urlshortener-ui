package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCMD = &cobra.Command{
	Use:     "version",
	Short:   "Shows version information",
	Example: "urlshortener-ui version",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s) Commit: %s by %s", Version, Date, Commit, BuiltBy)
	},
}

func init() {
	rootCmd.AddCommand(versionCMD)
}
