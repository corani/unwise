package mem

import (
	"sync"
	"time"

	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
)

type cachedBook struct {
	storage.Book

	mutex      sync.RWMutex
	Updated    time.Time
	Highlights []*storage.Highlight
}

type Mem struct {
	conf  *config.Config
	mutex sync.RWMutex
	books []*cachedBook
}

// Assert that Mem implements the storage.Storage interface.
var _ storage.Storage = (*Mem)(nil)

func New(conf *config.Config) *Mem {
	return &Mem{
		conf:  conf,
		mutex: sync.RWMutex{},
		books: nil,
	}
}

func (s *Mem) AddBook(title, author, source string) (storage.Book, bool) {
	if title == "" {
		title = storage.DefaultTitle
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, book := range s.books {
		if book.Title == title && book.Author == author && book.SourceURL == source {
			book.Book.Updated = book.Updated
			book.Book.Highlights = append([]*storage.Highlight{}, book.Highlights...)

			return book.Book, false
		}
	}

	book := &cachedBook{
		Book: storage.Book{
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

func (s *Mem) ListBooks(lt, gt time.Time) []storage.Book {
	if gt.IsZero() {
		gt = time.Now()
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	books := make([]storage.Book, 0, len(s.books))

	for _, book := range s.books {
		if book.Updated.Before(gt) && book.Updated.After(lt) {
			book.Book.Updated = book.Updated
			book.Book.Highlights = append([]*storage.Highlight{}, book.Highlights...)

			books = append(books, book.Book)
		}
	}

	return books
}

func (s *Mem) ListHighlights(lt, gt time.Time) []storage.Highlight {
	if gt.IsZero() {
		gt = time.Now()
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	highlights := make([]storage.Highlight, 0)

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

func (s *Mem) AddHighlight(b storage.Book, text, note, chapter string, location int, url string) (storage.Highlight, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, book := range s.books {
		if book.ID == b.ID {
			highlight, created := book.AddHighlight(text, note, chapter, location, url)

			return *highlight, created
		}
	}

	return storage.Highlight{}, false
}

func (b *cachedBook) AddHighlight(text, note, chapter string, location int, url string) (*storage.Highlight, bool) {
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

	highlight := &storage.Highlight{
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
