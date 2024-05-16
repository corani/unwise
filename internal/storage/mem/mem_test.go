package mem

import (
	"testing"
	"time"

	"github.com/corani/unwise/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMem_AddBook(t *testing.T) {
	rq := require.New(t)
	m := New(nil)

	b1, _ := m.AddBook("", "author", "source")
	rq.Len(m.books, 1)
	rq.Equal(storage.DefaultTitle, b1.Title)
	rq.Equal("author", b1.Author)
	rq.Equal("source", b1.SourceURL)

	// Add the same book again should update the existing
	// book.
	b2, _ := m.AddBook("", "author", "source")
	rq.Len(m.books, 1)
	rq.Equal(b1.ID, b2.ID)
	rq.Equal(b1.Title, b2.Title)
	rq.Equal(b1.Author, b2.Author)
	rq.Equal(b1.SourceURL, b2.SourceURL)
	rq.NotEqual(b1.Updated, b2.Updated)

	// Add a book with a different source.
	b3, _ := m.AddBook("", "author", "source2")
	rq.Len(m.books, 2)
	rq.NotEqual(b1.ID, b3.ID)
	rq.Equal(b1.Title, b3.Title)
	rq.Equal(b1.Author, b3.Author)
	rq.Equal("source2", b3.SourceURL)
	rq.NotEqual(b1.Updated, b2.Updated)
}

func TestMem_ListBooks(t *testing.T) {
	rq := require.New(t)
	m := New(nil)

	m.AddBook("title1", "author1", "source1")
	m.AddBook("title2", "author2", "source2")
	m.AddBook("title3", "author3", "source3")

	books := m.ListBooks(time.Time{}, time.Time{})
	rq.Len(books, 3)
}

func TestMem_AddHighlight(t *testing.T) {
	rq := require.New(t)
	m := New(nil)

	// New book
	b, created := m.AddBook("title", "author", "source")
	rq.True(created)

	// Add a highlight
	_, created = m.AddHighlight(b, "text1", "note1", "chapter1", 1, "url1")
	rq.True(created)

	// Add the same highlight again
	_, created = m.AddHighlight(b, "text1", "note3", "chapter3", 3, "url3")
	rq.False(created)

	// Add a different highlight
	_, created = m.AddHighlight(b, "text2", "note2", "chapter2", 2, "url2")
	rq.True(created)

	// Add a highlight for a non-existing book
	_, created = m.AddHighlight(storage.Book{ID: -1}, "text1", "note3", "chapter3", 3, "url3")
	rq.False(created)

	// Get the original book again
	b, created = m.AddBook("title", "author", "source")
	rq.False(created)
	rq.Len(b.Highlights, 2)
}

func TestMem_ListHighlights(t *testing.T) {
	rq := require.New(t)
	m := New(nil)

	// Add books
	b1, _ := m.AddBook("title1", "author", "source")
	b2, _ := m.AddBook("title2", "author", "source")

	// Add highlights
	m.AddHighlight(b1, "text1", "note1", "chapter1", 1, "url1")
	m.AddHighlight(b1, "text2", "note2", "chapter2", 2, "url2")
	m.AddHighlight(b2, "text3", "note3", "chapter3", 3, "url3")

	// List all highlights
	highlights := m.ListHighlights(time.Time{}, time.Time{})
	rq.Len(highlights, 3)
}
