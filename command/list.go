package command

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/k8sec/k8s"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	var (
		arguments    []string
		base64encode bool
		kubeconfig   string
		kubeClient   *client.Client
		namespace    string
		v            string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&base64encode, "base64", false, "If true, values are shown as base64-encoded string (Default: false)")
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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, strings.Join([]string{"NAME", "TYPE", "KEY", "VALUE"}, "\t"))

	if len(arguments) == 1 {
		secret, err := kubeClient.Secrets(namespace).Get(arguments[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		for key, value := range secret.Data {
			if base64encode {
				v = base64.StdEncoding.EncodeToString(value)
			} else {
				v = strconv.Quote(string(value))
			}

			fmt.Fprintln(w, strings.Join([]string{secret.Name, string(secret.Type), key, v}, "\t"))
		}
	} else {
		secrets, err := kubeClient.Secrets(namespace).List(api.ListOptions{})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		for _, secret := range secrets.Items {
			for key, value := range secret.Data {
				if base64encode {
					v = base64.StdEncoding.EncodeToString(value)
				} else {
					v = strconv.Quote(string(value))
				}

				fmt.Fprintln(w, strings.Join([]string{secret.Name, string(secret.Type), key, v}, "\t"))
			}
		}
	}

	w.Flush()

	return 0
}

func (c *ListCommand) Synopsis() string {
	return ""
}

func (c *ListCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
