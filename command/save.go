package command

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dtan4/k8sec/k8s"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type SaveCommand struct {
	Meta
}

func (c *SaveCommand) Run(args []string) int {
	var (
		arguments     []string
		base64encoded bool
		filename      string
		kubeconfig    string
		kubeClient    *client.Client
		namespace     string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&base64encoded, "base64", false, "If true, values are parsed as base64-encoded string (Default: false)")
	flags.StringVar(&filename, "f", "", "Path to save (Default: flush to stdout)")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file (Default: ~/.kube/config)")
	flags.StringVar(&namespace, "namespace", "", "Namespace scope (Default: default)")

	if err := flags.Parse(args[0:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	if len(arguments) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments.")
		return 1
	}

	if namespace == "" {
		namespace = api.NamespaceDefault
	}

	kubeClient, err := k8s.NewKubeClient(kubeconfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	var lines []string

	if len(arguments) == 1 {
		secret, err := kubeClient.Secrets(namespace).Get(arguments[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		for key, value := range secret.Data {
			lines = append(lines, key+"="+strconv.Quote(string(value)))
		}
	} else {
		secrets, err := kubeClient.Secrets(namespace).List(api.ListOptions{})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				lines = append(lines, key+"="+strconv.Quote(string(value)))
			}
		}
	}

	if filename != "" {
		f, err := os.Create(filename)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		defer f.Close()

		w := bufio.NewWriter(f)

		for _, line := range lines {
			_, err := w.WriteString(line + "\n")

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
		}

		w.Flush()
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	return 0
}

func (c *SaveCommand) Synopsis() string {
	return "Save as dotenv format"
}

func (c *SaveCommand) Help() string {
	helpText := `
$ k8sec save [--kubeconfig KUBECONFIG] [--namespace NAMESPACE] [-f FILENAME] [NAME]

Save as dotenv format

# Example
$ k8sec save rails
database-url="postgres://example.com:5432/dbname"

# Save as .env
$ k8sec save -f .env rails
$ cat .env
database-url="postgres://example.com:5432/dbname"
`
	return strings.TrimSpace(helpText)
}
