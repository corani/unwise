package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	defer func(old string) {
		version = old
	}(version)

	tt := []struct {
		set string
		exp string
	}{
		{
			set: " trim me ",
			exp: "trim me",
		},
		{
			set: "v0.0.3",
			exp: "v0.0.3",
		},
		{
			set: "16/merge",
			exp: "pr-16",
		},
	}

	for _, tc := range tt {
		t.Run(tc.exp, func(t *testing.T) {
			version = tc.set

			require.Equal(t, tc.exp, Version())
		})
	}
}

func TestHash(t *testing.T) {
	defer func(old string) {
		hash = old
	}(hash)

	tt := []struct {
		set string
		exp string
	}{
		{
			set: " trim me ",
			exp: "trim me",
		},
	}

	for _, tc := range tt {
		t.Run(tc.exp, func(t *testing.T) {
			hash = tc.set

			require.Equal(t, tc.exp, Hash())
		})
	}
}
