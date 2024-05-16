package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBook_NumHighlights(t *testing.T) {
	b := Book{
		Highlights: []*Highlight{
			{}, {}, {},
		},
	}

	require.Equal(t, 3, b.NumHighlights())
}

func TestBook_LastHighlight(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(1)
	t3 := t2.Add(1)

	b := Book{
		Highlights: []*Highlight{
			{Updated: t1}, {Updated: t2}, {Updated: t3},
		},
	}

	require.Equal(t, t3, b.LastHighlight())
}
