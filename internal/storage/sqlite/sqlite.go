package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
	_ "modernc.org/sqlite"
)

type DB struct {
	conf *config.Config
	db   *sql.DB
}

// Assert that Mem implements the storage.Storage interface.
var _ storage.Storage = (*DB)(nil)

func New(ctx context.Context, conf *config.Config) (*DB, error) {
	filename := filepath.Join(conf.DataPath, "quotes.db")

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%v?_pragma=journal_mode(WAL)", filename))
	if err != nil {
		return nil, err
	}

	s := &DB{
		conf: conf,
		db:   db,
	}

	return s, s.Init(ctx)
}

func (s *DB) Init(ctx context.Context) error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	if s.conf.DropTable == "true" {
		if _, err := s.db.ExecContext(ctx, `
			DROP TABLE IF EXISTS books;
			DROP TABLE IF EXISTS highlights;
		`); err != nil {
			return err
		}
	}

	if _, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS books (
			id 			INTEGER PRIMARY KEY AUTOINCREMENT,
			title 		TEXT NOT NULL DEFAULT '',
			author 		TEXT NOT NULL DEFAULT '',
			source_url  TEXT NOT NULL DEFAULT '',
			updated 	TEXT NOT NULL,
			UNIQUE (title, author, source_url)
		);
		CREATE TABLE IF NOT EXISTS highlights (
			id 			INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id 	INTEGER NOT NULL,
			text 		TEXT NOT NULL,
			note 		TEXT NOT NULL DEFAULT '',
			chapter 	TEXT NOT NULL DEFAULT '',
			location 	INTEGER NOT NULL DEFAULT 0,
			url 		TEXT NOT NULL DEFAULT '',
			updated 	TEXT NOT NULL,
			UNIQUE (book_id, text)
		);
	`); err != nil {
		return err
	}

	return nil
}

func (s *DB) AddBook(ctx context.Context, title, author, source string) (storage.Book, error) {
	if title == "" {
		title = storage.DefaultTitle
	}

	updated := time.Now()

	rows, err := s.db.QueryContext(ctx, `
		INSERT INTO books (title, author, source_url, updated) 
		VALUES (?, ?, ?, ?) 
		ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
		RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
	`, title, author, source, updated.Format(time.RFC3339), updated.Format(time.RFC3339))
	if err != nil {
		return storage.Book{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var id, numHighlights int

		if err := rows.Scan(&id, &numHighlights); err != nil {
			return storage.Book{}, err
		}

		return storage.Book{
			ID:            id,
			Title:         title,
			Author:        author,
			SourceURL:     source,
			Updated:       updated,
			NumHighlights: numHighlights,
		}, nil
	}

	return storage.Book{}, fmt.Errorf("no rows returned")
}

func (s *DB) ListBooks(ctx context.Context, lt, gt time.Time) ([]storage.Book, error) {
	if lt.IsZero() {
		lt = time.Now()
	}

	var books []storage.Book

	rows, err := s.db.QueryContext(ctx, `
		SELECT b.id, b.title, b.author, b.source_url, b.updated, COUNT(h.id) AS num_highlights
		FROM   books AS b LEFT OUTER JOIN highlights AS h ON b.id = h.book_id
		WHERE  b.updated >= ? AND b.updated <= ?
		GROUP BY b.id
	`, gt.Format(time.RFC3339), lt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			book    storage.Book
			updated string
		)

		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.SourceURL, &updated, &book.NumHighlights); err != nil {
			return nil, err
		}

		book.Updated, err = time.Parse(time.RFC3339, updated)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *DB) AddHighlight(ctx context.Context, b storage.Book, text, note, chapter string, location int, url string) (storage.Highlight, error) {
	updated := time.Now()

	rows, err := s.db.QueryContext(ctx, `
		INSERT INTO highlights (book_id, text, note, chapter, location, url, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?) 
		ON CONFLICT (book_id, text) DO UPDATE SET note = ?, chapter = ?, location = ?, url = ?, updated = ?
		RETURNING id
	`, b.ID, text, note, chapter, location, url, updated.Format(time.RFC3339),
		note, chapter, location, url, updated.Format(time.RFC3339))
	if err != nil {
		return storage.Highlight{}, err
	}

	defer rows.Close()

	if rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return storage.Highlight{}, err
		}

		return storage.Highlight{
			BookID:   b.ID,
			ID:       id,
			Text:     text,
			Note:     note,
			Chapter:  chapter,
			Location: location,
			URL:      url,
			Updated:  updated,
		}, nil
	}

	return storage.Highlight{}, fmt.Errorf("no rows returned")
}

func (s *DB) ListHighlights(ctx context.Context, lt, gt time.Time) ([]storage.Highlight, error) {
	if lt.IsZero() {
		lt = time.Now()
	}

	var highlights []storage.Highlight

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, book_id, text, note, chapter, location, url, updated 
		FROM   highlights 
		WHERE  updated >= ? AND updated <= ? 
	`, gt.Format(time.RFC3339), lt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			highlight storage.Highlight
			updated   string
		)

		if err := rows.Scan(&highlight.ID, &highlight.BookID, &highlight.Text, &highlight.Note, &highlight.Chapter, &highlight.Location, &highlight.URL, &updated); err != nil {
			return nil, err
		}

		highlight.Updated, err = time.Parse(time.RFC3339, updated)
		if err != nil {
			return nil, err
		}

		highlights = append(highlights, highlight)
	}

	return highlights, nil
}

func (s *DB) ListHighlightsByBook(ctx context.Context, bookID int) ([]storage.Highlight, error) {
	var highlights []storage.Highlight

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, book_id, text, note, chapter, location, url, updated 
		FROM   highlights 
		WHERE  book_id = ?
		ORDER BY location
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			highlight storage.Highlight
			updated   string
		)

		if err := rows.Scan(&highlight.ID, &highlight.BookID, &highlight.Text, &highlight.Note, &highlight.Chapter, &highlight.Location, &highlight.URL, &updated); err != nil {
			return nil, err
		}

		highlight.Updated, err = time.Parse(time.RFC3339, updated)
		if err != nil {
			return nil, err
		}

		highlights = append(highlights, highlight)
	}

	return highlights, nil
}
