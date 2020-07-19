package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/dtan4/k8sec/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var dumpOpts = struct {
	filename string
	noquotes bool
}{}

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump [NAME]",
	Short: "Dump secrets as dotenv (key=value) format",
	Long: `Dump secrets as dotenv (key=value) format

$ k8sec dump rails
database-url="postgres://example.com:5432/dbname"

Save as .env:

$ k8sec dump -f .env rails
$ cat .env
database-url="postgres://example.com:5432/dbname"

Save as .env without quotes:

$ k8sec dump -f .env --noquotes rails
$ cat .env
database-url=postgres://example.com:5432/dbname
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many arguments.")
		}

		ctx := context.Background()

		k8sclient, err := client.New(rootOpts.kubeconfig, rootOpts.context)
		if err != nil {
			return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
		}

		var namespace string

		if rootOpts.namespace != "" {
			namespace = rootOpts.namespace
		} else {
			namespace = k8sclient.DefaultNamespace()
		}

		return runDump(ctx, k8sclient, namespace, args, os.Stdout)
	},
}

func runDump(ctx context.Context, k8sclient client.Client, namespace string, args []string, out io.Writer) error {
	var lines []string

	if len(args) == 1 {
		secret, err := k8sclient.GetSecret(ctx, namespace, args[0])
		if err != nil {
			return errors.Wrapf(err, "Failed to get secret. name=%s", args[0])
		}

		for key, value := range secret.Data {
			line := string(value)
			if !dumpOpts.noquotes {
				line = strconv.Quote(line)
			}
			lines = append(lines, key+"="+line)
		}
	} else {
		secrets, err := k8sclient.ListSecrets(ctx, namespace)
		if err != nil {
			return errors.Wrap(err, "Failed to list secret.")
		}

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				v := string(value)
				if !dumpOpts.noquotes {
					v = strconv.Quote(v)
				}
				lines = append(lines, key+"="+v)
			}
		}
	}

	sort.Strings(lines)

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
			fmt.Fprintln(out, line)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringVarP(&dumpOpts.filename, "filename", "f", "", "File to dump")
	dumpCmd.Flags().BoolVar(&dumpOpts.noquotes, "noquotes", false, "Dump without quotes")
}
