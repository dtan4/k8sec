package cmd

import (
	"fmt"
	"io"

	"github.com/dtan4/k8sec/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(out io.Writer) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of k8sec",
		Run:   doVersion,
	}

	return versionCmd
}

func doVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version.String())
}
