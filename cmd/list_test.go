package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRunList(t *testing.T) {
	testcases := map[string]struct {
		base64encode bool
		args         []string
		secret       *v1.Secret
		secrets      *v1.SecretList
		err          error
		wantOut      string
		wantErr      error
	}{
		"no secret arg": {
			base64encode: false,
			args:         []string{},
			secrets: &v1.SecretList{
				Items: []v1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "default-token-12345",
						},
						Data: map[string][]byte{
							"ca.crt":    []byte("thisiscrt"),
							"namespace": []byte("test"),
							"token":     []byte("thisistoken"),
						},
						Type: v1.SecretTypeServiceAccountToken,
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "rails",
						},
						Data: map[string][]byte{
							"rails-env":    []byte("production"),
							"database-url": []byte("postgres://example.com:5432/dbname"),
						},
						Type: v1.SecretTypeOpaque,
					},
				},
			},
			err: nil,
			wantOut: `NAME			TYPE					KEY		VALUE
default-token-12345	kubernetes.io/service-account-token	ca.crt		"thisiscrt"
default-token-12345	kubernetes.io/service-account-token	namespace	"test"
default-token-12345	kubernetes.io/service-account-token	token		"thisistoken"
rails			Opaque					database-url	"postgres://example.com:5432/dbname"
rails			Opaque					rails-env	"production"
`,
			wantErr: nil,
		},

		"no secret arg and error": {
			args:    []string{},
			err:     errors.New("cannot retrieve secret rails"),
			wantErr: errors.New("list secrets: cannot retrieve secret rails"),
		},

		"one secret arg": {
			base64encode: false,
			args:         []string{"rails"},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rails",
				},
				Data: map[string][]byte{
					"rails-env":    []byte("production"),
					"database-url": []byte("postgres://example.com:5432/dbname"),
				},
				Type: v1.SecretTypeOpaque,
			},
			err: nil,
			wantOut: `NAME	TYPE	KEY		VALUE
rails	Opaque	database-url	"postgres://example.com:5432/dbname"
rails	Opaque	rails-env	"production"
`,
			wantErr: nil,
		},

		"one secret arg with --base64 option": {
			base64encode: true,
			args:         []string{"rails"},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rails",
				},
				Data: map[string][]byte{
					"rails-env":    []byte("production"),
					"database-url": []byte("postgres://example.com:5432/dbname"),
				},
				Type: v1.SecretTypeOpaque,
			},
			err: nil,
			wantOut: `NAME	TYPE	KEY		VALUE
rails	Opaque	database-url	cG9zdGdyZXM6Ly9leGFtcGxlLmNvbTo1NDMyL2RibmFtZQ==
rails	Opaque	rails-env	cHJvZHVjdGlvbg==
`,
			wantErr: nil,
		},

		"one secret and error": {
			args:    []string{"rails"},
			err:     errors.New("cannot retrieve secret rails"),
			wantErr: errors.New(`get secret "rails": cannot retrieve secret rails`),
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

			opts := listOpts{
				base64encode: tc.base64encode,
			}
			err := runList(context.Background(), k8sclient, namespace, tc.args, &out, &opts)

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
