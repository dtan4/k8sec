package cmd

import (
	"bytes"
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRunSet(t *testing.T) {
	testcases := map[string]struct {
		args    []string
		secret  *v1.Secret
		secrets *v1.SecretList
		err     error
		wantOut string
		wantErr error
	}{
		"create one key-value pair": {
			args: []string{
				"rails",
				"database-url=postgres://example.com:5432/dbname",
			},
			secrets: &v1.SecretList{
				Items: []v1.Secret{},
			},
			wantOut: "rails\n",
			wantErr: nil,
		},

		"update one key-value pair": {
			args: []string{
				"rails",
				"database-url=postgres://example.com:5432/dbname",
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rails",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			secrets: &v1.SecretList{
				Items: []v1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "rails",
						},
						Data: map[string][]byte{
							"foo": []byte("bar"),
						},
					},
				},
			},
			wantOut: "rails\n",
			wantErr: nil,
		},

		"two key-value pairs": {
			args: []string{
				"rails",
				"database-url=postgres://example.com:5432/dbname",
				"rails-env=production",
			},
			secrets: &v1.SecretList{
				Items: []v1.Secret{},
			},
			wantOut: "rails\n",
			wantErr: nil,
		},

		"error at get secret": {
			args: []string{
				"rails",
				"database-url=postgres://example.com:5432/dbname",
				"rails-env=production",
			},
			secrets: &v1.SecretList{
				Items: []v1.Secret{},
			},
			err:     fmt.Errorf("cannot list secret"),
			wantErr: fmt.Errorf("Failed to get current secret. name=rails: cannot list secret"),
		},
	}

	namespace := "test"

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			k8sclient := &fakeClient{
				getSecretResponse:   tc.secret,
				listSecretsResponse: tc.secrets,
				err:                 tc.err,
			}

			var out bytes.Buffer

			err := runSet(k8sclient, namespace, tc.args, &out)

			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("want error %q, got no error", tc.wantErr.Error())
				}

				if err.Error() != tc.wantErr.Error() {
					t.Fatalf("want error %q, got %q", tc.wantErr.Error(), err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, got %q", err.Error())
				}

				if out.String() != tc.wantOut {
					t.Logf("want:\n%s", tc.wantOut)
					t.Logf("got:\n%s", out.String())
					t.Fatalf("want %q, got %q", tc.wantOut, out.String())
				}
			}
		})
	}
}
