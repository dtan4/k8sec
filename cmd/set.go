package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dtan4/k8sec/pkg/client"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

type setOpts struct {
	base64encoded bool
}

func newSetCmd(out io.Writer) *cobra.Command {
	opts := setOpts{}

	setCmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("too few arguments")
			}

			ctx := context.Background()

			k8sclient, err := client.New(rootOpts.kubeconfig, rootOpts.context)
			if err != nil {
				return fmt.Errorf("initialize Kubernetes API client: %w", err)
			}

			var namespace string

			if rootOpts.namespace != "" {
				namespace = rootOpts.namespace
			} else {
				namespace = k8sclient.DefaultNamespace()
			}

			return runSet(ctx, k8sclient, namespace, args, out, &opts)
		},
	}

	setCmd.Flags().BoolVar(&opts.base64encoded, "base64", false, "Decode the given value as base64-encoded string")

	return setCmd
}

func runSet(ctx context.Context, k8sclient client.Client, namespace string, args []string, out io.Writer, opts *setOpts) error {
	name := args[0]

	data := map[string][]byte{}

	for _, kv := range args[1:] {
		ary := strings.SplitN(kv, "=", 2)

		if len(ary) != 2 {
			return errors.New("argument should be in key=value format argument")
		}

		k, v := ary[0], ary[1]

		if opts.base64encoded {
			_v, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return fmt.Errorf("decode value as base64-encoded string: %w", err)
			}

			data[k] = _v
		} else {
			data[k] = []byte(v)
		}
	}

	ss, err := k8sclient.ListSecrets(ctx, namespace)
	if err != nil {
		return fmt.Errorf("get current secret %q: %w", name, err)
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
		s, err = k8sclient.GetSecret(ctx, namespace, name)
		if err != nil {
			return fmt.Errorf("get current secret %q: %w", name, err)
		}

		if s.Data == nil {
			s.Data = data
		} else {
			for k, v := range data {
				s.Data[k] = v
			}
		}

		_, err = k8sclient.UpdateSecret(ctx, namespace, s)
		if err != nil {
			return fmt.Errorf("update secret %q: %w", name, err)
		}
	} else {
		s = &v1.Secret{
			Data: data,
		}
		s.SetName(name)

		_, err = k8sclient.CreateSecret(ctx, namespace, s)
		if err != nil {
			return fmt.Errorf("create secret %q: %w", name, err)
		}
	}

	fmt.Fprintln(out, s.Name)

	return nil
}
