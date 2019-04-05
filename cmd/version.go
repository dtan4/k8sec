package cmd

import (
	"fmt"

	"github.com/wantedly/k8sec/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of k8sec",
	Run:   doVersion,
}

func doVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version.String())
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
