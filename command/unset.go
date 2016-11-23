package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/k8sec/k8s"
	"k8s.io/client-go/pkg/api/v1"
)

type UnsetCommand struct {
	Meta
}

func (c *UnsetCommand) Run(args []string) int {
	var (
		arguments  []string
		kubeconfig string
		namespace  string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

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

	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, "Too few arguments. Example: $ k8sec unset rails RAILS_ENV")
		return 1
	}

	name := arguments[0]

	clientset, err := k8s.NewKubeClient(kubeconfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	s, err := clientset.Core().Secrets(namespace).Get(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for _, k := range arguments[1:] {
		_, ok := s.Data[k]
		if !ok {
			fmt.Fprintln(os.Stderr, "The key "+k+" does not exist.")
			return 1
		}

		delete(s.Data, k)
	}

	_, err = clientset.Core().Secrets(namespace).Update(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Println(s.Name)

	return 0
}

func (c *UnsetCommand) Synopsis() string {
	return "Unset secrets"
}

func (c *UnsetCommand) Help() string {
	helpText := `
$ k8sec unset [--kubeconfig KUBECONFIG] [--namespace NAMESPACE] NAME KEY1 KEY2

Unset secrets

# Example
$ k8sec unset rails rails-env
`
	return strings.TrimSpace(helpText)
}
