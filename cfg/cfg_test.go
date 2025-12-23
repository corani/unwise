package cfg

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func setFS(t *testing.T, files map[string]string) {
	t.Helper()

	mapfs := make(fstest.MapFS, len(files))
	for name, content := range files {
		mapfs[name] = &fstest.MapFile{Data: []byte(content)}
	}

	oldFS := getFS
	getFS = func() fs.FS {
		return mapfs
	}

	t.Cleanup(func() {
		getFS = oldFS
	})
}

func TestVersion(t *testing.T) {
	tt := []struct {
		name  string
		files map[string]string
		exp   string
	}{
		{
			name:  "simple",
			files: map[string]string{"VERSION": "v1.2.3"},
			exp:   "v1.2.3",
		},
		{
			name:  "with trim",
			files: map[string]string{"VERSION": "  v1.2.4  "},
			exp:   "v1.2.4",
		},
		{
			name:  "pr",
			files: map[string]string{"VERSION": "123/merge"},
			exp:   "pr-123",
		},
		{
			name:  "pr with trim",
			files: map[string]string{"VERSION": " 124/merge "},
			exp:   "pr-124",
		},
		{
			name:  "empty",
			files: map[string]string{"VERSION": ""},
			exp:   "",
		},
		{
			name: "no file",
			exp:  "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			setFS(t, tc.files)

			require.Equal(t, tc.exp, Version())
		})
	}
}

func TestHash(t *testing.T) {
	tt := []struct {
		name  string
		files map[string]string
		exp   string
	}{
		{
			name:  "simple",
			files: map[string]string{"HASH": "some-hash"},
			exp:   "some-hash",
		},
		{
			name:  "with trim",
			files: map[string]string{"HASH": "  some-other-hash  "},
			exp:   "some-other-hash",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			setFS(t, tc.files)
			require.Equal(t, tc.exp, Hash())
		})
	}
}

func TestBuild(t *testing.T) {
	tt := []struct {
		name  string
		files map[string]string
		exp   string
	}{
		{
			name:  "simple",
			files: map[string]string{"BUILD": "some-build-date"},
			exp:   "some-build-date",
		},
		{
			name:  "with trim",
			files: map[string]string{"BUILD": "  some-other-build-date  "},
			exp:   "some-other-build-date",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			setFS(t, tc.files)
			require.Equal(t, tc.exp, Build())
		})
	}
}
