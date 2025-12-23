package cfg

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed *
var files embed.FS

// getFS returns the filesystem to read configuration files from. This is needed to adapt the
// type to an `fs.FS` so we can inject a mock during testing.
//
//nolint:gochecknoglobals
var getFS = func() fs.FS {
	return files
}

func getStringFile(name string) string {
	bs, err := fs.ReadFile(getFS(), name)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(bs))
}

// Version returns the embedded application version.
func Version() string {
	v := getStringFile("VERSION")

	if before, ok := strings.CutSuffix(v, "/merge"); ok {
		v = "pr-" + before
	}

	return v
}

// Hash returns the embedded application hash.
func Hash() string {
	return getStringFile("HASH")
}

// Build returns the embedded application hash.
func Build() string {
	return getStringFile("BUILD")
}
