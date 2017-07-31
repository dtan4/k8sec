package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootOpts = struct {
	context    string
	debug      bool
	kubeconfig string
	namespace  string
}{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "k8sec",
	Short:         "CLI tool to manage Kubernetes Secrets easily",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if rootOpts.debug {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&rootOpts.context, "context", "", "Kubernetes context")
	RootCmd.PersistentFlags().BoolVar(&rootOpts.debug, "debug", false, "Debug mode")
	RootCmd.PersistentFlags().StringVar(&rootOpts.kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	RootCmd.PersistentFlags().StringVarP(&rootOpts.namespace, "namespace", "n", "", "Kubernetes namespace")
}

func initConfig() {
}
