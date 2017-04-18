package version

import (
	"testing"
)

func TestString(t *testing.T) {
	Version = "v0.1.0"
	Revision = "abcd1234"

	expected := "k8sec version v0.1.0, build abcd1234"
	actual := String()

	if actual != expected {
		t.Errorf("Version string does not match. expected: %q, actual: %q", expected, actual)
	}
}
