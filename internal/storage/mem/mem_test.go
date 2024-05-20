package mem

import (
	"context"
	"testing"
	"time"

	"github.com/corani/unwise/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMem_AddBook(t *testing.T) {
	rq := require.New(t)
	m := New(nil)

	b1, err := m.AddBook(context.Background(), "", "author", "source")
	rq.NoError(err)
	rq.Len(m.books, 1)
	rq.Equal(storage.DefaultTitle, b1.Title)
	rq.Equal("author", b1.Author)
	rq.Equal("source", b1.SourceURL)

	// Add the same book again should update the existing
	// book.
	b2, _ := m.AddBook(context.Background(), "", "author", "source")
	rq.Len(m.books, 1)
	rq.Equal(b1.ID, b2.ID)
	rq.Equal(b1.Title, b2.Title)
	rq.Equal(b1.Author, b2.Author)
	rq.Equal(b1.SourceURL, b2.SourceURL)
	rq.NotEqual(b1.Updated, b2.Updated)

	// Add a book with a different source.
	b3, _ := m.AddBook(context.Background(), "", "author", "source2")
	rq.Len(m.books, 2)
	rq.NotEqual(b1.ID, b3.ID)
	rq.Equal(b1.Title, b3.Title)
	rq.Equal(b1.Author, b3.Author)
	rq.Equal("source2", b3.SourceURL)
	rq.NotEqual(b1.Updated, b2.Updated)
}

func TestMem_ListBooks(t *testing.T) {
	ctx := context.Background()
	rq := require.New(t)
	m := New(nil)

	_, _ = m.AddBook(ctx, "title1", "author1", "source1")
	_, _ = m.AddBook(ctx, "title2", "author2", "source2")
	_, _ = m.AddBook(ctx, "title3", "author3", "source3")

	books, err := m.ListBooks(ctx, time.Time{}, time.Time{})
	rq.NoError(err)
	rq.Len(books, 3)
}

func TestMem_AddHighlight(t *testing.T) {
	rq := require.New(t)
	m := New(nil)
	ctx := context.Background()

	// New book
	b, err := m.AddBook(ctx, "title", "author", "source")
	rq.NoError(err)

	// Add a highlight
	_, err = m.AddHighlight(ctx, b, "text1", "note1", "chapter1", 1, "url1")
	rq.NoError(err)

	// Add the same highlight again
	_, err = m.AddHighlight(ctx, b, "text1", "note3", "chapter3", 3, "url3")
	rq.NoError(err)

	// Add a different highlight
	_, err = m.AddHighlight(ctx, b, "text2", "note2", "chapter2", 2, "url2")
	rq.NoError(err)

	// Add a highlight for a non-existing book
	_, err = m.AddHighlight(ctx, storage.Book{ID: -1}, "text1", "note3", "chapter3", 3, "url3")
	rq.NoError(err)

	// Get the original book again
	b, err = m.AddBook(ctx, "title", "author", "source")
	rq.NoError(err)
	rq.Equal(2, b.NumHighlights)
}

func TestMem_ListHighlights(t *testing.T) {
	ctx := context.Background()
	rq := require.New(t)
	m := New(nil)

	// Add books
	b1, _ := m.AddBook(ctx, "title1", "author", "source")
	b2, _ := m.AddBook(ctx, "title2", "author", "source")

	// Add highlights
	_, _ = m.AddHighlight(ctx, b1, "text1", "note1", "chapter1", 1, "url1")
	_, _ = m.AddHighlight(ctx, b1, "text2", "note2", "chapter2", 2, "url2")
	_, _ = m.AddHighlight(ctx, b2, "text3", "note3", "chapter3", 3, "url3")

	// List all highlights
	highlights, err := m.ListHighlights(ctx, time.Time{}, time.Time{})
	rq.NoError(err)
	rq.Len(highlights, 3)
}
