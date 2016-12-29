package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/dtan4/k8sec/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api/v1"
)

var dumpOpts = struct {
	filename string
}{}

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: doDump,
}

func doDump(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Too many arguments.")
	}

	clientset, err := k8s.NewKubeClient(rootOpts.kubeconfig)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
	}

	var lines []string

	if len(args) == 1 {
		secret, err := clientset.Core().Secrets(rootOpts.namespace).Get(args[0])
		if err != nil {
			return errors.Wrapf(err, "Failed to get secret. name=%s", args[0])
		}

		for key, value := range secret.Data {
			lines = append(lines, key+"="+strconv.Quote(string(value)))
		}
	} else {
		secrets, err := clientset.Core().Secrets(rootOpts.namespace).List(v1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "Failed to list secret.")
		}

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				lines = append(lines, key+"="+strconv.Quote(string(value)))
			}
		}
	}

	if dumpOpts.filename != "" {
		f, err := os.Create(dumpOpts.filename)
		if err != nil {
			return errors.Wrapf(err, "Failed to open file. filename=%s", dumpOpts.filename)
		}
		defer f.Close()

		w := bufio.NewWriter(f)

		for _, line := range lines {
			_, err := w.WriteString(line + "\n")
			if err != nil {
				return errors.Wrapf(err, "Failed to write to file. filename=%s", dumpOpts.filename)
			}
		}

		w.Flush()
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringVarP(&dumpOpts.filename, "filename", "f", "", "File to dump")
}
