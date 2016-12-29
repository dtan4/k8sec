package cmd

import (
	"fmt"

	"github.com/dtan4/k8sec/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doVersion,
}

func doVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version.String())
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
