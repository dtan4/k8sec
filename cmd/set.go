package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/dtan4/k8sec/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var setOpts = struct {
	base64encoded bool
}{}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: doSet,
}

func doSet(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Too few arguments.")
	}
	name := args[0]

	data := map[string][]byte{}

	for _, kv := range args[1:] {
		ary := strings.SplitN(kv, "=", 2)

		if len(args) != 2 {
			return errors.Errorf("Argument should be in key=value format. argument=%q", kv)
		}

		k, v := ary[0], ary[1]

		if setOpts.base64encoded {
			_v, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return errors.Wrapf(err, "Failed to decode value as base64-encoded string. value=%q", v)
			}

			data[k] = _v
		} else {
			data[k] = []byte(v)
		}
	}

	clientset, err := k8s.NewKubeClient(rootOpts.kubeconfig)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
	}

	s, err := clientset.Core().Secrets(rootOpts.namespace).Get(name)
	if err != nil {
		return errors.Wrapf(err, "Failed to get current secret. name=%s", name)
	}

	for k, v := range data {
		s.Data[k] = v
	}

	_, err = clientset.Core().Secrets(rootOpts.namespace).Update(s)
	if err != nil {
		return errors.Wrapf(err, "Failed to set secret. name=%s", name)
	}

	fmt.Println(s.Name)

	return nil
}

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.Flags().BoolVar(&setOpts.base64encoded, "base64", false, "Decode the given value as base64-encoded string")
}
