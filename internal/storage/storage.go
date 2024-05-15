package storage

import (
	"sync"
	"time"

	"github.com/corani/unwise/internal/config"
)

const DefaultTitle = "Quotes"

type CachedBook struct {
	Book

	mutex      sync.RWMutex
	Updated    time.Time
	Highlights []*Highlight
}

type Book struct {
	ID        int
	Title     string
	Author    string
	SourceURL string

	// cloned from CachedBook
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

type Storage struct {
	conf  *config.Config
	mutex sync.RWMutex
	books []*CachedBook
}

func New(conf *config.Config) *Storage {
	return &Storage{
		conf:  conf,
		mutex: sync.RWMutex{},
		books: nil,
	}
}

// The highlights array can be length 1+ and each highlight can be from the same or multiple
// books/articles. If you don't include a title, we'll put the highlight in a generic "Quotes"
// book, and if you don't include an author we'll keep it blank or just use the URL domain (if
// a source_url was provided).
//
// Finally, we de-dupe highlights by title/author/text/source_url. So if you send a highlight
// with those 4 things the same (including nulls) then it will do nothing rather than create a
// "duplicate".

func (s *Storage) AddBook(title, author, source string) (Book, bool) {
	if title == "" {
		title = DefaultTitle
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, book := range s.books {
		if book.Title == title && book.Author == author && book.SourceURL == source {
			book.Book.Updated = book.Updated
			book.Book.Highlights = append([]*Highlight{}, book.Highlights...)

			return book.Book, false
		}
	}

	book := &CachedBook{
		Book: Book{
			ID:        len(s.books) + 1,
			Title:     title,
			Author:    author,
			SourceURL: source,
			Updated:   time.Now(),
		},
		Updated: time.Now(),
	}

	s.books = append(s.books, book)

	return book.Book, true
}

func (s *Storage) ListBooks(lt, gt time.Time) []Book {
	if gt.IsZero() {
		gt = time.Now()
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	books := make([]Book, 0, len(s.books))

	for _, book := range s.books {
		if book.Updated.Before(gt) && book.Updated.After(lt) {
			book.Book.Updated = book.Updated
			book.Book.Highlights = append([]*Highlight{}, book.Highlights...)

			books = append(books, book.Book)
		}
	}

	return books
}

func (s *Storage) ListHighlights(lt, gt time.Time) []Highlight {
	if gt.IsZero() {
		gt = time.Now()
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	highlights := make([]Highlight, 0)

	for _, book := range s.books {
		book.mutex.RLock()
		for _, highlight := range book.Highlights {
			if highlight.Updated.Before(gt) && highlight.Updated.After(lt) {
				highlights = append(highlights, *highlight)
			}
		}
		book.mutex.RUnlock()
	}

	return highlights
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

func (s *Storage) AddHighlight(b Book, text, note, chapter string, location int, url string) (Highlight, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, book := range s.books {
		if book.ID == b.ID {
			highlight, created := book.AddHighlight(text, note, chapter, location, url)

			return *highlight, created
		}
	}

	return Highlight{}, false
}

func (b *CachedBook) AddHighlight(text, note, chapter string, location int, url string) (*Highlight, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.Updated = time.Now()

	// if a highlight with the same text exists, update it. Otherwise add a new one.
	for _, h := range b.Highlights {
		if h.Text == text {
			h.Note = note
			h.Chapter = chapter
			h.Location = location
			h.URL = url
			h.Updated = b.Updated

			return h, false
		}
	}

	highlight := &Highlight{
		ID:       len(b.Highlights) + 1,
		BookID:   b.ID,
		Text:     text,
		Note:     note,
		Chapter:  chapter,
		Location: location,
		URL:      url,
		Updated:  b.Updated,
	}

	b.Highlights = append(b.Highlights, highlight)

	return highlight, true
}
