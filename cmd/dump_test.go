package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRunDump(t *testing.T) {
	testcases := map[string]struct {
		args     []string
		filename string
		secret   *v1.Secret
		secrets  *v1.SecretList
		err      error
		wantOut  string
		wantErr  error
	}{
		"no secret arg": {
			args:     []string{},
			filename: "",
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
			wantOut: `ca.crt="thisiscrt"
database-url="postgres://example.com:5432/dbname"
namespace="test"
rails-env="production"
token="thisistoken"
`,
			wantErr: nil,
		},

		"no secret arg and error": {
			args:     []string{},
			filename: "",
			err:      fmt.Errorf("cannot retrieve secret rails"),
			wantErr:  fmt.Errorf("Failed to list secret.: cannot retrieve secret rails"),
		},

		"one secret arg": {
			args:     []string{"rails"},
			filename: "",
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
			wantOut: `database-url="postgres://example.com:5432/dbname"
rails-env="production"
`,
			wantErr: nil,
		},

		// TODO: Add testcase for --noquotes once I move noquotes to local variable

		"one secret and error": {
			args:     []string{"rails"},
			filename: "",
			err:      fmt.Errorf("cannot retrieve secret rails"),
			wantErr:  fmt.Errorf("Failed to get secret. name=rails: cannot retrieve secret rails"),
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

			opts := dumpOpts{}

			err := runDump(context.Background(), k8sclient, namespace, tc.args, &out, &opts)

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

func TestRunDump_dumpToFile(t *testing.T) {
	testcases := map[string]struct {
		args     []string
		filename string
		secret   *v1.Secret
		secrets  *v1.SecretList
		wantOut  string
		wantBody string
	}{
		"dump to file": {
			args:     []string{},
			filename: ".env",
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
			wantOut: "",
			wantBody: `ca.crt="thisiscrt"
database-url="postgres://example.com:5432/dbname"
namespace="test"
rails-env="production"
token="thisistoken"
`,
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
			}

			var out bytes.Buffer

			filename := filepath.Join(t.TempDir(), tc.filename)

			opts := dumpOpts{
				filename: filename,
			}

			err := runDump(context.Background(), k8sclient, namespace, tc.args, &out, &opts)

			if err != nil {
				t.Fatalf("want no error, got %q", err.Error())
			}

			if out.String() != tc.wantOut {
				t.Logf("want:\n%s", tc.wantOut)
				t.Logf("got:\n%s", out.String())
				t.Fatalf("want %q, got %q", tc.wantOut, out.String())
			}

			if _, err := os.Stat(filename); err != nil {
				t.Fatalf("want file %q but not found", filename)
			}

			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}

			b, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			if string(b) != tc.wantBody {
				t.Fatalf("want %q, got %q", tc.wantBody, string(b))
			}
		})
	}
}
