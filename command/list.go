package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/fields"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	var (
		arguments     []string
		kubeconfig    string
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

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	if kubeconfig == "" {
		loadingRules.ExplicitPath = clientcmd.RecommendedHomeFile
	} else {
		loadingRules.ExplicitPath = kubeconfig
	}

	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	clientConfig, err := loader.ClientConfig()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	kubeclient, err := client.New(clientConfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	secrets, err := kubeclient.Secrets(api.NamespaceDefault).List(api.ListOptions{
		FieldSelector: fieldSelector,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, strings.Join([]string{"NAMESPACE", "NAME", "KEY", "VALUE"}, "\t"))

	for _, secret := range secrets.Items {
		for key, value := range secret.Data {
			fmt.Fprintln(w, strings.Join([]string{secret.Namespace, secret.Name, key, strconv.Quote(string(value))}, "\t"))
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
