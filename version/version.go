package version

import (
	"fmt"
)

var (
	// Version represents version number
	Version string
	// Revision represents commit hash at built binary
	Revision string
)

// String returns version string
func String() string {
	return fmt.Sprintf("k8sec version %s, build %s", Version, Revision)
}
