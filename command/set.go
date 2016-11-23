package command

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/k8sec/k8s"
	"k8s.io/client-go/pkg/api/v1"
)

type SetCommand struct {
	Meta
}

func (c *SetCommand) Run(args []string) int {
	var (
		arguments     []string
		base64encoded bool
		kubeconfig    string
		namespace     string
	)

	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&base64encoded, "base64", false, "If true, values are parsed as base64-encoded string (Default: false)")
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
		fmt.Fprintln(os.Stderr, "Too few arguments. Example: $ k8sec set rails RAILS_ENV=production")
		return 1
	}

	name := arguments[0]

	data := map[string][]byte{}

	for _, kv := range arguments[1:] {
		ary := strings.SplitN(kv, "=", 2)

		if len(ary) != 2 {
			fmt.Fprintln(os.Stderr, "Argument should be key=value format. Given argument: "+kv)
			return 1
		}

		k, v := ary[0], ary[1]

		if base64encoded {
			_v, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}

			data[k] = _v
		} else {
			data[k] = []byte(v)
		}
	}

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

	for k, v := range data {
		s.Data[k] = v
	}

	_, err = clientset.Core().Secrets(namespace).Update(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Println(s.Name)

	return 0
}

func (c *SetCommand) Synopsis() string {
	return "Set secrets"
}

func (c *SetCommand) Help() string {
	helpText := `
$ k8sec set [--base64] [--kubeconfig KUBECONFIG] [--namespace NAMESPACE] NAME KEY1=VALUE1 KEY2=VALUE2

Set secrets

# Example
$ k8sec set rails rails-env=production
rails

# Pass base64-encoded value
$ echo dtan4 | base64
ZHRhbjQK
$ k8sec set --base64 rails foo=ZHRhbjQK
rails
$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"
rails   Opaque  foo             "dtan4\n"
`
	return strings.TrimSpace(helpText)
}
