package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/k8sec/k8s"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	var (
		arguments     []string
		kubeconfig    string
		kubeClient    *client.Client
		fieldSelector fields.Selector
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")

	if err := flags.Parse(args[0:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	if len(arguments) >= 1 {
		fieldSelector = fields.Set{api.ObjectNameField: arguments[0]}.AsSelector()
	} else {
		fieldSelector = fields.Everything()
	}

	kubeClient, err := k8s.NewKubeClient(kubeconfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	secrets, err := kubeClient.Secrets(api.NamespaceDefault).List(api.ListOptions{
		FieldSelector: fieldSelector,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, strings.Join([]string{"NAME", "TYPE", "KEY", "VALUE"}, "\t"))

	for _, secret := range secrets.Items {
		for key, value := range secret.Data {
			fmt.Fprintln(w, strings.Join([]string{secret.Name, string(secret.Type), key, strconv.Quote(string(value))}, "\t"))
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
