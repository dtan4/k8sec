package command

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dtan4/k8sec/k8s"
	"k8s.io/client-go/pkg/api/v1"
)

type SaveCommand struct {
	Meta
}

func (c *SaveCommand) Run(args []string) int {
	var (
		arguments  []string
		filename   string
		kubeconfig string
		namespace  string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.StringVar(&filename, "f", "", "Path to save (Default: flush to stdout)")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file (Default: ~/.kube/config)")
	flags.StringVar(&namespace, "namespace", v1.NamespaceDefault, "Namespace scope")

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

	clientset, err := k8s.NewKubeClient(kubeconfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	var lines []string

	if len(arguments) == 1 {
		secret, err := clientset.Core().Secrets(namespace).Get(arguments[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		for key, value := range secret.Data {
			lines = append(lines, key+"="+strconv.Quote(string(value)))
		}
	} else {
		secrets, err := clientset.Core().Secrets(namespace).List(v1.ListOptions{})
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
	return "Save as dotenv (key=value) format"
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
