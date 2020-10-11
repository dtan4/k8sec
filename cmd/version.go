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
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(args, out)
		},
	}

	return versionCmd
}

func runVersion(args []string, out io.Writer) {
	fmt.Fprintln(out, version.String())
}
