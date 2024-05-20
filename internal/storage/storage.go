package storage

import (
	"context"
	"time"
)

const DefaultTitle = "Quotes"

type Storage interface {
	// AddBook adds a new book to the storage. If the book already exists, it returns the
	// existing book.
	AddBook(ctx context.Context, title, author, source string) (Book, error)

	// AddHighlight adds a new highlight to the storage. If the highlight already exists,
	// it returns the existing highlight.
	AddHighlight(ctx context.Context, b Book, text, note, chapter string, location int, url string) (Highlight, error)

	// ListBooks returns a list of books from the storage.
	ListBooks(ctx context.Context, lt, gt time.Time) ([]Book, error)

	// ListHighlights returns a list of highlights from the storage.
	ListHighlights(ctx context.Context, lt, gt time.Time) ([]Highlight, error)
}

type Book struct {
	ID            int
	Title         string
	Author        string
	SourceURL     string
	Updated       time.Time
	NumHighlights int
	LastHighlight time.Time
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
