package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dtan4/k8sec/pkg/client"
	"github.com/spf13/cobra"
)

type loadOpts struct {
	filename string
}

func newLoadCmd(in io.Reader, out io.Writer) *cobra.Command {
	opts := loadOpts{}

	loadCmd := &cobra.Command{
		Use:   "load NAME",
		Short: "Load secrets from dotenv (key=value) format text",
		Long: `Load secrets from dotenv (key=value) format text

$ cat .env
database-url="postgres://example.com:5432/dbname"
$ k8sec load -f .env rails

Load from stdin:

$ cat .env | k8sec load rails
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("variable name must be specified.")
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

			return runLoad(ctx, k8sclient, namespace, args, in, out, &opts)
		},
	}

	loadCmd.Flags().StringVarP(&opts.filename, "filename", "f", "", "File to load")

	return loadCmd
}

func runLoad(ctx context.Context, k8sclient client.Client, namespace string, args []string, in io.Reader, out io.Writer, opts *loadOpts) error {
	if len(args) != 1 {
		return fmt.Errorf("Variable name must be specified.")
	}
	name := args[0]

	var sc *bufio.Scanner
	data := map[string][]byte{}

	if opts.filename != "" {
		f, err := os.Open(opts.filename)
		if err != nil {
			return fmt.Errorf("open file %q: %w", opts.filename, err)
		}
		defer f.Close()

		sc = bufio.NewScanner(f)
	} else {
		sc = bufio.NewScanner(in)
	}

	for sc.Scan() {
		line := sc.Text()
		ary := strings.SplitN(line, "=", 2)

		if len(ary) != 2 {
			return errors.New("line must be key=value format")
		}

		k, v := ary[0], ary[1]

		_v, err := strconv.Unquote(v)
		if err != nil {
			// Parse as is
			_v = v
		}

		data[k] = []byte(_v)
	}

	s, err := k8sclient.GetSecret(ctx, namespace, name)
	if err != nil {
		return fmt.Errorf("get secret %q: %w", name, err)
	}

	for k, v := range data {
		s.Data[k] = v
	}

	_, err = k8sclient.UpdateSecret(ctx, namespace, s)
	if err != nil {
		return fmt.Errorf("set secret %q: %w", name, err)
	}

	return nil
}
