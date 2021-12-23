package version

import (
	"testing"
)

func TestString(t *testing.T) {
	version = "v0.1.0"
	commit = "751282c0a11e8bacd2d1d9597728c6cb3f9147f8"
	date = "2020-03-12T16:34:36Z"

	want := "k8sec version: v0.1.0, commit: 751282c0a11e8bacd2d1d9597728c6cb3f9147f8, build at: 2020-03-12T16:34:36Z"
	got := String()

	if got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
}

func TestUserAgent(t *testing.T) {
	version = "v0.1.0"
	commit = "751282c0a11e8bacd2d1d9597728c6cb3f9147f8"
	date = "2020-03-12T16:34:36Z"

	want := "k8sec/v0.1.0"
	got := UserAgent()

	if got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
}
