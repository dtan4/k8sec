package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRunLoad(t *testing.T) {
	testcases := map[string]struct {
		args    []string
		secret  *v1.Secret
		input   string
		err     error
		wantErr error
	}{
		// TODO: Add appropriate error handling test for "no secret"

		"update one key-value pair": {
			args: []string{
				"rails",
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rails",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			input:   `database-url="postgres://example.com:5432/dbname"`,
			wantErr: nil,
		},

		"two key-value pairs": {
			args: []string{
				"rails",
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rails",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			input: `database-url="postgres://example.com:5432/dbname"
rails-env="production"`,
			wantErr: nil,
		},

		"error at get secret": {
			args: []string{
				"rails",
			},
			err:     fmt.Errorf("cannot list secret"),
			wantErr: fmt.Errorf("Failed to get secret. name=rails: cannot list secret"),
		},
	}

	namespace := "test"

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			k8sclient := &fakeClient{
				getSecretResponse: tc.secret,
				err:               tc.err,
			}

			in := strings.NewReader(tc.input)
			var out bytes.Buffer

			err := runLoad(k8sclient, namespace, tc.args, in, &out)

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
			}
		})
	}
}
