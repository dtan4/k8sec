package version

import (
	"fmt"
)

var (
	version string
	commit  string
	date    string
)

// String returns version string
func String() string {
	return fmt.Sprintf("k8sec version: %s, commit: %s, build at: %s", version, commit, date)
}

// UserAgent returns the version string in User-Agent header style
func UserAgent() string {
	return fmt.Sprintf("k8sec/%s", version)
}
