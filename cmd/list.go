package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/k8sec/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type listOpts struct {
	base64encode bool
}

func newListCmd(out io.Writer) *cobra.Command {
	opts := listOpts{}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets",
		Long: `List secrets

$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"

Show values as base64-encoded string:

$ k8sec list --base64 rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    cG9zdGdyZXM6Ly9leGFtcGxlLmNvbTo1NDMyL2RibmFtZQ==
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("Too many arguments.")
			}

			ctx := context.Background()

			k8sclient, err := client.New(rootOpts.kubeconfig, rootOpts.context)
			if err != nil {
				return errors.Wrap(err, "Failed to initialize Kubernetes API client.")
			}

			var namespace string

			if rootOpts.namespace != "" {
				namespace = rootOpts.namespace
			} else {
				namespace = k8sclient.DefaultNamespace()
			}

			return runList(ctx, k8sclient, namespace, args, out, &opts)
		},
	}

	listCmd.Flags().BoolVar(&opts.base64encode, "base64", false, "Show values as base64-encoded string")

	return listCmd
}

type Secret struct {
	Name  string
	Type  string
	Key   string
	Value string
}

func runList(ctx context.Context, k8sclient client.Client, namespace string, args []string, out io.Writer, opts *listOpts) error {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, strings.Join([]string{"NAME", "TYPE", "KEY", "VALUE"}, "\t"))

	var v string

	secrets := []Secret{}

	if len(args) == 1 {
		secret, err := k8sclient.GetSecret(ctx, namespace, args[0])
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve secrets.")
		}

		for key, value := range secret.Data {
			if opts.base64encode {
				v = base64.StdEncoding.EncodeToString(value)
			} else {
				v = strconv.Quote(string(value))
			}

			secrets = append(secrets, Secret{
				Name:  secret.Name,
				Type:  string(secret.Type),
				Key:   key,
				Value: v,
			})
		}

		// sort by KEY
		sort.Slice(secrets, func(i, j int) bool {
			return secrets[i].Key < secrets[j].Key
		})
	} else {
		ss, err := k8sclient.ListSecrets(ctx, namespace)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve secrets.")
		}

		for _, secret := range ss.Items {
			kvs := []struct {
				k, v string
			}{}

			for key, value := range secret.Data {
				if opts.base64encode {
					v = base64.StdEncoding.EncodeToString(value)
				} else {
					v = strconv.Quote(string(value))
				}

				kvs = append(kvs, struct {
					k, v string
				}{
					k: key,
					v: v,
				})
			}

			// sort by KEY
			sort.Slice(kvs, func(i, j int) bool {
				return kvs[i].k < kvs[j].k
			})

			for _, kv := range kvs {
				secrets = append(secrets, Secret{
					Name:  secret.Name,
					Type:  string(secret.Type),
					Key:   kv.k,
					Value: kv.v,
				})
			}
		}
	}

	for _, secret := range secrets {
		fmt.Fprintln(w, strings.Join([]string{secret.Name, secret.Type, secret.Key, secret.Value}, "\t"))
	}

	w.Flush()

	return nil
}
