package storage

import (
	"time"
)

const DefaultTitle = "Quotes"

type Storage interface {
	// AddBook adds a new book to the storage. If the book already exists, it returns the
	// existing book.
	AddBook(title, author, source string) (Book, bool)

	// AddHighlight adds a new highlight to the storage. If the highlight already exists,
	// it returns the existing highlight.
	AddHighlight(b Book, text, note, chapter string, location int, url string) (Highlight, bool)

	// ListBooks returns a list of books from the storage.
	ListBooks(lt, gt time.Time) []Book

	// ListHighlights returns a list of highlights from the storage.
	ListHighlights(lt, gt time.Time) []Highlight
}

type Book struct {
	ID         int
	Title      string
	Author     string
	SourceURL  string
	Updated    time.Time
	Highlights []*Highlight
}

type Highlight struct {
	BookID   int
	ID       int
	Text     string
	Note     string
	Chapter  string
	Location int
	URL      string
	Updated  time.Time
}

func (b *Book) NumHighlights() int {
	return len(b.Highlights)
}

func (b *Book) LastHighlight() time.Time {
	last := time.Time{}

	for _, h := range b.Highlights {
		if h.Updated.After(last) {
			last = h.Updated
		}
	}

	return last
}
