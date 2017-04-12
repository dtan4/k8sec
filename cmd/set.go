package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/dtan4/k8sec/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api/v1"
)

var setOpts = struct {
	base64encoded bool
}{}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set NAME KEY1=VALUE1 [KEY2=VALUE2 ...]",
	Short: "Set secrets",
	Long: `Set secrets

Set value as it is:

$ k8sec set rails rails-env=production
rails

Set base64-encoded value:

$ echo -n dtan4 | base64
ZHRhbjQ=
$ k8sec set --base64 rails foo=ZHRhbjQ=
rails

Result:

$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"
rails   Opaque  foo             "dtan4"
`,
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

		if len(ary) != 2 {
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

	k8sclient, err := k8s.NewKubeClient(rootOpts.kubeconfig, rootOpts.context, rootOpts.namespace)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
	}

	ss, err := k8sclient.ListSecrets()
	if err != nil {
		return errors.Wrapf(err, "Failed to get current secret. name=%s", name)
	}

	exists := false

	for _, s := range ss.Items {
		if s.Name == name {
			exists = true
			break
		}
	}

	var s *v1.Secret

	if exists {
		s, err = k8sclient.Clientset.Core().Secrets(k8sclient.Namespace).Get(name)
		if err != nil {
			return errors.Wrapf(err, "Failed to get current secret. name=%s", name)
		}

		if s.Data == nil {
			s.Data = data
		} else {
			for k, v := range data {
				s.Data[k] = v
			}
		}

		_, err = k8sclient.Clientset.Core().Secrets(k8sclient.Namespace).Update(s)
		if err != nil {
			return errors.Wrapf(err, "Failed to update secret. name=%s", name)
		}
	} else {
		s = &v1.Secret{
			Data: data,
		}
		s.SetName(name)

		_, err = k8sclient.Clientset.Core().Secrets(k8sclient.Namespace).Create(s)
		if err != nil {
			return errors.Wrapf(err, "Failed to create secret. name=%s", name)
		}
	}

	fmt.Println(s.Name)

	return nil
}

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.Flags().BoolVar(&setOpts.base64encoded, "base64", false, "Decode the given value as base64-encoded string")
}
