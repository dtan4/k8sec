package cmd

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dtan4/k8sec/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var loadOpts = struct {
	filename string
}{}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
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
			return errors.New("Variable name must be specified.")
		}

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

		return runLoad(k8sclient, namespace, args, os.Stdin, os.Stdout)
	},
}

func runLoad(k8sclient client.Client, namespace string, args []string, in io.Reader, out io.Writer) error {
	if len(args) != 1 {
		return errors.New("Variable name must be specified.")
	}
	name := args[0]

	var sc *bufio.Scanner
	data := map[string][]byte{}

	if loadOpts.filename != "" {
		f, err := os.Open(loadOpts.filename)
		if err != nil {
			return errors.Wrapf(err, "Failed to open file. filename=%s", loadOpts.filename)
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
			return errors.Errorf("Line should be key=value format. line=%q", line)
		}

		k, v := ary[0], ary[1]

		_v, err := strconv.Unquote(v)
		if err != nil {
			// Parse as is
			_v = v
		}

		data[k] = []byte(_v)
	}

	s, err := k8sclient.GetSecret(namespace, name)
	if err != nil {
		return errors.Wrapf(err, "Failed to get secret. name=%s", name)
	}

	for k, v := range data {
		s.Data[k] = v
	}

	_, err = k8sclient.UpdateSecret(namespace, s)
	if err != nil {
		return errors.Wrapf(err, "Failed to set secret. name=%s", name)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(loadCmd)

	loadCmd.Flags().StringVarP(&loadOpts.filename, "filename", "f", "", "File to load")
}
