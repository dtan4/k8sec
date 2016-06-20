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
		kubeconfig    string
		kubeClient    *client.Client
		namespace     string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&base64encoded, "base64", false, "If true, values are parsed as base64-encoded string (Default: false)")
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

	secrets, err := kubeClient.Secrets(namespace).List(api.ListOptions{})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if len(arguments) == 1 {
		f, err := os.Create(arguments[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		defer f.Close()

		w := bufio.NewWriter(f)

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				_, err := w.WriteString(key + "=" + strconv.Quote(string(value)) + "\n")

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return 1
				}
			}
		}

		w.Flush()
	} else {
		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				fmt.Println(key + "=" + strconv.Quote(string(value)))
			}
		}
	}

	return 0
}

func (c *SaveCommand) Synopsis() string {
	return "Save as dotenv format"
}

func (c *SaveCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
