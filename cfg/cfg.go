package cfg

import (
	_ "embed" // needed for `go:embed`
	"strings"
)

var (
	//go:embed VERSION
	version string

	//go:embed HASH
	hash string

	//go:embed BUILD
	build string
)

// Version returns the embedded application version.
func Version() string {
	v := strings.TrimSpace(version)

	if strings.HasSuffix(v, "/merge") {
		v = "pr-" + strings.TrimSuffix(v, "/merge")
	}

	return v
}

// Hash returns the embedded application hash.
func Hash() string {
	return strings.TrimSpace(hash)
}

// Build returns the embedded application hash.
func Build() string {
	return strings.TrimSpace(build)
}
