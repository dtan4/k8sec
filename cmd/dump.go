package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/dtan4/k8sec/pkg/client"
	"github.com/spf13/cobra"
)

type dumpOpts struct {
	filename string
	noquotes bool
}

func newDumpCmd(out io.Writer) *cobra.Command {
	opts := dumpOpts{}

	dumpCmd := &cobra.Command{
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
				return errors.New("too many arguments")
			}

			ctx := context.Background()

			k8sclient, err := client.New(rootOpts.kubeconfig, rootOpts.context)
			if err != nil {
				return fmt.Errorf("initialize Kubernetes API client: %w", err)
			}

			var namespace string

			if rootOpts.namespace != "" {
				namespace = rootOpts.namespace
			} else {
				namespace = k8sclient.DefaultNamespace()
			}

			return runDump(ctx, k8sclient, namespace, args, out, &opts)
		},
	}

	dumpCmd.Flags().StringVarP(&opts.filename, "filename", "f", "", "File to dump")
	dumpCmd.Flags().BoolVar(&opts.noquotes, "noquotes", false, "Dump without quotes")

	return dumpCmd
}

func runDump(ctx context.Context, k8sclient client.Client, namespace string, args []string, out io.Writer, opts *dumpOpts) error {
	var lines []string

	if len(args) == 1 {
		secret, err := k8sclient.GetSecret(ctx, namespace, args[0])
		if err != nil {
			return fmt.Errorf("get secret %q: %w", args[0], err)
		}

		for key, value := range secret.Data {
			line := string(value)
			if !opts.noquotes {
				line = strconv.Quote(line)
			}
			lines = append(lines, key+"="+line)
		}
	} else {
		secrets, err := k8sclient.ListSecrets(ctx, namespace)
		if err != nil {
			return fmt.Errorf("list secret: %w", err)
		}

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				v := string(value)
				if !opts.noquotes {
					v = strconv.Quote(v)
				}
				lines = append(lines, key+"="+v)
			}
		}
	}

	sort.Strings(lines)

	if opts.filename != "" {
		f, err := os.Create(opts.filename)
		if err != nil {
			return fmt.Errorf("open file %q: %w", opts.filename, err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)

		for _, line := range lines {
			_, err := w.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("write to file %q: %w", opts.filename, err)
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
