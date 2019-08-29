package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/dtan4/k8sec/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var dumpOpts = struct {
	filename string
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
`,
	RunE: doDump,
}

func doDump(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Too many arguments.")
	}

	k8sclient, err := k8s.NewKubeClient(rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
	}

	var namespace string

	if rootOpts.namespace != "" {
		namespace = rootOpts.namespace
	} else {
		namespace = k8sclient.DefaultNamespace()
	}

	var lines []string

	if len(args) == 1 {
		secret, err := k8sclient.GetSecret(namespace, args[0])
		if err != nil {
			return errors.Wrapf(err, "Failed to get secret. name=%s", args[0])
		}

		for key, value := range secret.Data {
			lines = append(lines, key+"="+strconv.Quote(string(value)))
		}
	} else {
		secrets, err := k8sclient.ListSecrets(namespace)
		if err != nil {
			return errors.Wrap(err, "Failed to list secret.")
		}

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				lines = append(lines, key+"="+strconv.Quote(string(value)))
			}
		}
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})

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
