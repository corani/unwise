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
)

// Version returns the embedded application version.
func Version() string {
	return strings.TrimSpace(version)
}

// Hash returns the embedded application hash.
func Hash() string {
	return strings.TrimSpace(hash)
}
