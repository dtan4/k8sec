package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var rootOpts = struct {
	context    string
	debug      bool
	kubeconfig string
	namespace  string
}{}

func newRootCmd(in io.Reader, out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           "k8sec",
		Short:         "CLI tool to manage Kubernetes Secrets easily",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	flags := cmd.PersistentFlags()

	flags.StringVar(&rootOpts.context, "context", "", "Kubernetes context")
	flags.BoolVar(&rootOpts.debug, "debug", false, "Debug mode")
	flags.StringVar(&rootOpts.kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	flags.StringVarP(&rootOpts.namespace, "namespace", "n", "", "Kubernetes namespace")

	cmd.AddCommand(newDumpCmd(out))
	cmd.AddCommand(newListCmd(out))
	cmd.AddCommand(newLoadCmd(in, out))
	cmd.AddCommand(newSetCmd(out))
	cmd.AddCommand(newUnsetCmd(out))
	cmd.AddCommand(newVersionCmd(out))

	return cmd
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(in io.Reader, out io.Writer, args []string) {
	cmd := newRootCmd(in, out, args)

	if err := cmd.Execute(); err != nil {
		if rootOpts.debug {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	}
}
