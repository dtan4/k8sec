package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUnSet(t *testing.T) {
	testcases := map[string]struct {
		args    []string
		secret  *v1.Secret
		err     error
		wantOut string
		wantErr error
	}{
		"delete secret": {
			args: []string{
				"rails",
				"database-url",
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rails",
				},
				Data: map[string][]byte{
					"database-url": []byte("postgres://example.com:5432/dbname"),
				},
			},
			wantOut: "rails\n",
			wantErr: nil,
		},

		// TODO: add testcase for no matched error found

		"error at get secret": {
			args: []string{
				"rails",
				"database-url",
			},
			err:     errors.New("cannot list secret"),
			wantErr: errors.New(`get current secret "rails": cannot list secret`),
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

			var out bytes.Buffer

			err := runUnset(context.Background(), k8sclient, namespace, tc.args, &out)

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
