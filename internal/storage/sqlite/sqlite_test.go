package sqlite

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func queryMatcher(t *testing.T) sqlmock.QueryMatcherFunc {
	t.Helper()

	return sqlmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		space := regexp.MustCompile(`\s+`)

		expectedSQL = space.ReplaceAllString(expectedSQL, " ")
		actualSQL = space.ReplaceAllString(actualSQL, " ")

		if expectedSQL != actualSQL {
			t.Logf("expectedSQL: %v", expectedSQL)
			t.Logf("actualSQL: %v", actualSQL)

			return fmt.Errorf("not equal")
		}

		return nil
	})
}

func TestSqlite_Init(t *testing.T) {
	conf := config.MustLoad()
	conf.DropTable = "true"

	tt := []struct {
		name   string
		setup  func(*testing.T, sqlmock.Sqlmock)
		expErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectPing()
				mock.ExpectExec(`
					DROP TABLE IF EXISTS books; 
					DROP TABLE IF EXISTS highlights;
				`).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(`
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
				`).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expErr: false,
		},
		{
			name: "ping error",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name: "error 1",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectPing()
				mock.ExpectExec(`
					DROP TABLE IF EXISTS books; 
					DROP TABLE IF EXISTS highlights;
				`).WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name: "error 2",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectPing()
				mock.ExpectExec(`
					DROP TABLE IF EXISTS books; 
					DROP TABLE IF EXISTS highlights;
				`).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(`
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
				`).WillReturnError(assert.AnError)
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.MonitorPingsOption(true),
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				conf: conf,
				db:   db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				rq.Error(s.Init(context.Background()))
			} else {
				rq.NoError(s.Init(context.Background()))
			}
		})
	}
}

func TestSqlite_AddBook(t *testing.T) {
	conf := config.MustLoad()

	tt := []struct {
		name    string
		title   string
		author  string
		source  string
		setup   func(*testing.T, sqlmock.Sqlmock, string)
		expErr  bool
		expBook storage.Book
	}{
		{
			name:   "success",
			title:  "title",
			author: "author",
			source: "source",
			setup: func(t *testing.T, mock sqlmock.Sqlmock, expTitle string) {
				mock.ExpectQuery(`
					INSERT INTO books (title, author, source_url, updated) 
					VALUES (?, ?, ?, ?) 
					ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
					RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
				`).WithArgs(expTitle, "author", "source", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "count"}).AddRow(1, 2))
			},
			expErr: false,
			expBook: storage.Book{
				ID:            1,
				Title:         "title",
				Author:        "author",
				SourceURL:     "source",
				NumHighlights: 2,
			},
		},
		{
			name:   "default title",
			title:  "",
			author: "author",
			source: "source",
			setup: func(t *testing.T, mock sqlmock.Sqlmock, expTitle string) {
				mock.ExpectQuery(`
					INSERT INTO books (title, author, source_url, updated) 
					VALUES (?, ?, ?, ?) 
					ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
					RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
				`).WithArgs(expTitle, "author", "source", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "count"}).AddRow(1, 2))
			},
			expErr: false,
			expBook: storage.Book{
				ID:            1,
				Title:         storage.DefaultTitle,
				Author:        "author",
				SourceURL:     "source",
				NumHighlights: 2,
			},
		},
		{
			name:   "query error",
			title:  "title",
			author: "author",
			source: "source",
			setup: func(t *testing.T, mock sqlmock.Sqlmock, expTitle string) {
				mock.ExpectQuery(`
					INSERT INTO books (title, author, source_url, updated) 
					VALUES (?, ?, ?, ?) 
					ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
					RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
				`).WithArgs(expTitle, "author", "source", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name:   "scan error",
			title:  "title",
			author: "author",
			source: "source",
			setup: func(t *testing.T, mock sqlmock.Sqlmock, expTitle string) {
				mock.ExpectQuery(`
					INSERT INTO books (title, author, source_url, updated) 
					VALUES (?, ?, ?, ?) 
					ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
					RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
				`).WithArgs(expTitle, "author", "source", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "count"}).AddRow("invalid", "invalid"))
			},
			expErr: true,
		},
		{
			name:   "no rows",
			title:  "title",
			author: "author",
			source: "source",
			setup: func(t *testing.T, mock sqlmock.Sqlmock, expTitle string) {
				mock.ExpectQuery(`
					INSERT INTO books (title, author, source_url, updated) 
					VALUES (?, ?, ?, ?) 
					ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
					RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
				`).WithArgs(expTitle, "author", "source", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "count"}))
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.MonitorPingsOption(true),
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				conf: conf,
				db:   db,
			}

			expTitle := tc.title
			if expTitle == "" {
				expTitle = storage.DefaultTitle
			}
			tc.setup(t, mock, expTitle)

			if tc.expErr {
				_, err := s.AddBook(context.Background(), tc.title, tc.author, tc.source)
				rq.Error(err)
			} else {
				book, err := s.AddBook(context.Background(), tc.title, tc.author, tc.source)
				rq.NoError(err)

				rq.Equal(tc.expBook.ID, book.ID)
				rq.Equal(tc.expBook.Title, book.Title)
				rq.Equal(tc.expBook.Author, book.Author)
				rq.Equal(tc.expBook.SourceURL, book.SourceURL)
				rq.Equal(tc.expBook.NumHighlights, book.NumHighlights)
			}
		})
	}
}

func TestSqlite_ListBooks(t *testing.T) {
	now := time.Now()

	tt := []struct {
		name     string
		lt       time.Time
		gt       time.Time
		setup    func(*testing.T, sqlmock.Sqlmock)
		expErr   bool
		expBooks []storage.Book
	}{
		{
			name: "success",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT b.id, b.title, b.author, b.source_url, b.updated, COUNT(h.id) AS num_highlights
					FROM   books AS b LEFT OUTER JOIN highlights AS h ON b.id = h.book_id
					WHERE  b.updated >= ? AND b.updated <= ?
					GROUP BY b.id
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "source_url", "updated", "num_highlights"}).
						AddRow(1, "title", "author", "source", now, 2))
			},
			expErr: false,
			expBooks: []storage.Book{
				{
					ID:            1,
					Title:         "title",
					Author:        "author",
					SourceURL:     "source",
					NumHighlights: 2,
				},
			},
		},
		{
			name: "query error",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT b.id, b.title, b.author, b.source_url, b.updated, COUNT(h.id) AS num_highlights
					FROM   books AS b LEFT OUTER JOIN highlights AS h ON b.id = h.book_id
					WHERE  b.updated >= ? AND b.updated <= ?
					GROUP BY b.id
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name: "scan error",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT b.id, b.title, b.author, b.source_url, b.updated, COUNT(h.id) AS num_highlights
					FROM   books AS b LEFT OUTER JOIN highlights AS h ON b.id = h.book_id
					WHERE  b.updated >= ? AND b.updated <= ?
					GROUP BY b.id
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "source_url", "updated", "num_highlights"}).
						AddRow("invalid", "title", "author", "source", now, 2))
			},
			expErr: true,
		},
		{
			name: "time error",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT b.id, b.title, b.author, b.source_url, b.updated, COUNT(h.id) AS num_highlights
					FROM   books AS b LEFT OUTER JOIN highlights AS h ON b.id = h.book_id
					WHERE  b.updated >= ? AND b.updated <= ?
					GROUP BY b.id
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "source_url", "updated", "num_highlights"}).
						AddRow(1, "title", "author", "source", "invalid", 2))
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.MonitorPingsOption(true),
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				db: db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				_, err := s.ListBooks(context.Background(), tc.lt, tc.gt)
				rq.Error(err)
			} else {
				books, err := s.ListBooks(context.Background(), tc.lt, tc.gt)
				rq.NoError(err)

				rq.Equal(len(tc.expBooks), len(books))

				for i, expBook := range tc.expBooks {
					rq.Equal(expBook.ID, books[i].ID)
					rq.Equal(expBook.Title, books[i].Title)
					rq.Equal(expBook.Author, books[i].Author)
					rq.Equal(expBook.SourceURL, books[i].SourceURL)
					rq.Equal(expBook.NumHighlights, books[i].NumHighlights)
				}
			}
		})
	}
}

func TestSqlite_AddHighlight(t *testing.T) {
	now := time.Now()

	tt := []struct {
		name         string
		book         storage.Book
		text         string
		note         string
		chapter      string
		location     int
		url          string
		setup        func(*testing.T, sqlmock.Sqlmock)
		expErr       bool
		expHighlight storage.Highlight
	}{
		{
			name: "success",
			book: storage.Book{
				ID: 1,
			},
			text:     "text",
			note:     "note",
			chapter:  "chapter",
			location: 3,
			url:      "url",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					INSERT INTO highlights (book_id, text, note, chapter, location, url, updated)
					VALUES (?, ?, ?, ?, ?, ?, ?) 
					ON CONFLICT (book_id, text) DO UPDATE SET note = ?, chapter = ?, location = ?, url = ?, updated = ?
					RETURNING id
				`).WithArgs(
					1, "text", "note", "chapter", 3, "url", now.Format(time.RFC3339),
					"note", "chapter", 3, "url", now.Format(time.RFC3339),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expErr: false,
			expHighlight: storage.Highlight{
				BookID:   1,
				ID:       1,
				Text:     "text",
				Note:     "note",
				Chapter:  "chapter",
				Location: 3,
				URL:      "url",
			},
		},
		{
			name: "query error",
			book: storage.Book{
				ID: 1,
			},
			text:     "text",
			note:     "note",
			chapter:  "chapter",
			location: 3,
			url:      "url",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					INSERT INTO highlights (book_id, text, note, chapter, location, url, updated)
					VALUES (?, ?, ?, ?, ?, ?, ?) 
					ON CONFLICT (book_id, text) DO UPDATE SET note = ?, chapter = ?, location = ?, url = ?, updated = ?
					RETURNING id
				`).WithArgs(
					1, "text", "note", "chapter", 3, "url", now.Format(time.RFC3339),
					"note", "chapter", 3, "url", now.Format(time.RFC3339),
				).WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name: "scan error",
			book: storage.Book{
				ID: 1,
			},
			text:     "text",
			note:     "note",
			chapter:  "chapter",
			location: 3,
			url:      "url",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					INSERT INTO highlights (book_id, text, note, chapter, location, url, updated)
					VALUES (?, ?, ?, ?, ?, ?, ?) 
					ON CONFLICT (book_id, text) DO UPDATE SET note = ?, chapter = ?, location = ?, url = ?, updated = ?
					RETURNING id
				`).WithArgs(
					1, "text", "note", "chapter", 3, "url", now.Format(time.RFC3339),
					"note", "chapter", 3, "url", now.Format(time.RFC3339),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("invalid"))
			},
			expErr: true,
		},
		{
			name: "no rows",
			book: storage.Book{
				ID: 1,
			},
			text:     "text",
			note:     "note",
			chapter:  "chapter",
			location: 3,
			url:      "url",
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					INSERT INTO highlights (book_id, text, note, chapter, location, url, updated)
					VALUES (?, ?, ?, ?, ?, ?, ?) 
					ON CONFLICT (book_id, text) DO UPDATE SET note = ?, chapter = ?, location = ?, url = ?, updated = ?
					RETURNING id
				`).WithArgs(
					1, "text", "note", "chapter", 3, "url", now.Format(time.RFC3339),
					"note", "chapter", 3, "url", now.Format(time.RFC3339),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.MonitorPingsOption(true),
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				db: db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				_, err := s.AddHighlight(context.Background(), tc.book, tc.text, tc.note, tc.chapter, tc.location, tc.url)
				rq.Error(err)
			} else {
				highlight, err := s.AddHighlight(context.Background(), tc.book, tc.text, tc.note, tc.chapter, tc.location, tc.url)
				rq.NoError(err)

				rq.Equal(tc.expHighlight.BookID, highlight.BookID)
				rq.Equal(tc.expHighlight.ID, highlight.ID)
				rq.Equal(tc.expHighlight.Text, highlight.Text)
				rq.Equal(tc.expHighlight.Note, highlight.Note)
				rq.Equal(tc.expHighlight.Chapter, highlight.Chapter)
				rq.Equal(tc.expHighlight.Location, highlight.Location)
				rq.Equal(tc.expHighlight.URL, highlight.URL)
			}
		})
	}
}

func TestSqlite_ListHighlightsByBook(t *testing.T) {
	now := time.Now()

	tt := []struct {
		name          string
		bookID        int
		setup         func(*testing.T, sqlmock.Sqlmock)
		expErr        bool
		expHighlights []storage.Highlight
	}{
		{
			name:   "success",
			bookID: 1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM highlights 
					WHERE book_id = ? 
					ORDER BY location
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
						AddRow(1, 1, "text1", "note1", "chapter1", 10, "url1", now).
						AddRow(2, 1, "text2", "note2", "chapter2", 20, "url2", now))
			},
			expErr: false,
			expHighlights: []storage.Highlight{
				{
					BookID:   1,
					ID:       1,
					Text:     "text1",
					Note:     "note1",
					Chapter:  "chapter1",
					Location: 10,
					URL:      "url1",
					Updated:  now,
				},
				{
					BookID:   1,
					ID:       2,
					Text:     "text2",
					Note:     "note2",
					Chapter:  "chapter2",
					Location: 20,
					URL:      "url2",
					Updated:  now,
				},
			},
		},
		{
			name:   "query error",
			bookID: 1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM highlights 
					WHERE book_id = ? 
					ORDER BY location
				`).WithArgs(1).
					WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name:   "scan error",
			bookID: 1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM highlights 
					WHERE book_id = ? 
					ORDER BY location
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
						AddRow("invalid", 1, "text", "note", "chapter", 3, "url", now))
			},
			expErr: true,
		},
		{
			name:   "time error",
			bookID: 1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM highlights 
					WHERE book_id = ? 
					ORDER BY location
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
						AddRow(1, 1, "text", "note", "chapter", 3, "url", "invalid"))
			},
			expErr: true,
		},
		{
			name:   "no highlights",
			bookID: 999,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM highlights 
					WHERE book_id = ? 
					ORDER BY location
				`).WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}))
			},
			expErr:        false,
			expHighlights: []storage.Highlight{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.MonitorPingsOption(true),
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				db: db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				_, err := s.ListHighlightsByBook(context.Background(), tc.bookID)
				rq.Error(err)
			} else {
				highlights, err := s.ListHighlightsByBook(context.Background(), tc.bookID)
				rq.NoError(err)

				rq.Equal(len(tc.expHighlights), len(highlights))

				for i, expHighlight := range tc.expHighlights {
					rq.Equal(expHighlight.BookID, highlights[i].BookID)
					rq.Equal(expHighlight.ID, highlights[i].ID)
					rq.Equal(expHighlight.Text, highlights[i].Text)
					rq.Equal(expHighlight.Note, highlights[i].Note)
					rq.Equal(expHighlight.Chapter, highlights[i].Chapter)
					rq.Equal(expHighlight.Location, highlights[i].Location)
					rq.Equal(expHighlight.URL, highlights[i].URL)
				}
			}
		})
	}
}

func TestSqlite_ListHighlights(t *testing.T) {
	now := time.Now()

	tt := []struct {
		name          string
		lt            time.Time
		gt            time.Time
		setup         func(*testing.T, sqlmock.Sqlmock)
		expErr        bool
		expHighlights []storage.Highlight
	}{
		{
			name: "success",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM   highlights 
					WHERE  updated >= ? AND updated <= ? 
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
						AddRow(1, 1, "text", "note", "chapter", 3, "url", now))
			},
			expErr: false,
			expHighlights: []storage.Highlight{
				{
					BookID:   1,
					ID:       1,
					Text:     "text",
					Note:     "note",
					Chapter:  "chapter",
					Location: 3,
					URL:      "url",
					Updated:  now,
				},
			},
		},
		{
			name: "query error",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM   highlights 
					WHERE  updated >= ? AND updated <= ? 
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnError(assert.AnError)
			},
			expErr: true,
		},
		{
			name: "scan error",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM   highlights 
					WHERE  updated >= ? AND updated <= ? 
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
						AddRow("invalid", 1, "text", "note", "chapter", 3, "url", now))
			},
			expErr: true,
		},
		{
			name: "time error",
			lt:   now,
			gt:   now,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`
					SELECT id, book_id, text, note, chapter, location, url, updated 
					FROM   highlights 
					WHERE  updated >= ? AND updated <= ? 
				`).WithArgs(now.Format(time.RFC3339), now.Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
						AddRow(1, 1, "text", "note", "chapter", 3, "url", "invalid"))
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.MonitorPingsOption(true),
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				db: db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				_, err := s.ListHighlights(context.Background(), tc.lt, tc.gt)
				rq.Error(err)
			} else {
				highlights, err := s.ListHighlights(context.Background(), tc.lt, tc.gt)
				rq.NoError(err)

				rq.Equal(len(tc.expHighlights), len(highlights))

				for i, expBook := range tc.expHighlights {
					rq.Equal(expBook.BookID, highlights[i].BookID)
					rq.Equal(expBook.ID, highlights[i].ID)
					rq.Equal(expBook.Text, highlights[i].Text)
					rq.Equal(expBook.Note, highlights[i].Note)
					rq.Equal(expBook.Chapter, highlights[i].Chapter)
					rq.Equal(expBook.Location, highlights[i].Location)
					rq.Equal(expBook.URL, highlights[i].URL)
				}
			}
		})
	}
}

func TestSqlite_UpdateHighlight(t *testing.T) {
	tt := []struct {
		name         string
		highlight    storage.Highlight
		setup        func(*testing.T, sqlmock.Sqlmock)
		expErr       bool
		expHighlight storage.Highlight
	}{
		{
			name: "success",
			highlight: storage.Highlight{
				ID:       1,
				Text:     "updated text",
				Note:     "updated note",
				Chapter:  "updated chapter",
				Location: 5,
			},
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`
					UPDATE highlights 
					SET text = ?, note = ?, chapter = ?, location = ?, updated = ?
					WHERE id = ?
				`).WithArgs("updated text", "updated note", "updated chapter", 5, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(`
					SELECT book_id, url 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id", "url"}).AddRow(2, "original-url"))

				mock.ExpectExec(`
					UPDATE books 
					SET updated = ?
					WHERE id = ?
				`).WithArgs(sqlmock.AnyArg(), 2).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			expErr: false,
			expHighlight: storage.Highlight{
				ID:       1,
				BookID:   2,
				Text:     "updated text",
				Note:     "updated note",
				Chapter:  "updated chapter",
				Location: 5,
				URL:      "original-url",
			},
		},
		{
			name: "update error",
			highlight: storage.Highlight{
				ID:       1,
				Text:     "updated text",
				Note:     "updated note",
				Chapter:  "updated chapter",
				Location: 5,
			},
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`
					UPDATE highlights 
					SET text = ?, note = ?, chapter = ?, location = ?, updated = ?
					WHERE id = ?
				`).WithArgs("updated text", "updated note", "updated chapter", 5, sqlmock.AnyArg(), 1).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expErr: true,
		},
		{
			name: "query error",
			highlight: storage.Highlight{
				ID:       1,
				Text:     "updated text",
				Note:     "updated note",
				Chapter:  "updated chapter",
				Location: 5,
			},
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`
					UPDATE highlights 
					SET text = ?, note = ?, chapter = ?, location = ?, updated = ?
					WHERE id = ?
				`).WithArgs("updated text", "updated note", "updated chapter", 5, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(`
					SELECT book_id, url 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expErr: true,
		},
		{
			name: "scan error",
			highlight: storage.Highlight{
				ID:       1,
				Text:     "updated text",
				Note:     "updated note",
				Chapter:  "updated chapter",
				Location: 5,
			},
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`
					UPDATE highlights 
					SET text = ?, note = ?, chapter = ?, location = ?, updated = ?
					WHERE id = ?
				`).WithArgs("updated text", "updated note", "updated chapter", 5, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(`
					SELECT book_id, url 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id", "url"}).AddRow("invalid", "url"))
				mock.ExpectRollback()
			},
			expErr: true,
		},
		{
			name: "book update error",
			highlight: storage.Highlight{
				ID:       1,
				Text:     "updated text",
				Note:     "updated note",
				Chapter:  "updated chapter",
				Location: 5,
			},
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`
					UPDATE highlights 
					SET text = ?, note = ?, chapter = ?, location = ?, updated = ?
					WHERE id = ?
				`).WithArgs("updated text", "updated note", "updated chapter", 5, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(`
					SELECT book_id, url 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id", "url"}).AddRow(2, "original-url"))

				mock.ExpectExec(`
					UPDATE books 
					SET updated = ?
					WHERE id = ?
				`).WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				db: db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				_, err := s.UpdateHighlight(context.Background(), tc.highlight)
				rq.Error(err)
			} else {
				highlight, err := s.UpdateHighlight(context.Background(), tc.highlight)
				rq.NoError(err)

				rq.Equal(tc.expHighlight.ID, highlight.ID)
				rq.Equal(tc.expHighlight.BookID, highlight.BookID)
				rq.Equal(tc.expHighlight.Text, highlight.Text)
				rq.Equal(tc.expHighlight.Note, highlight.Note)
				rq.Equal(tc.expHighlight.Chapter, highlight.Chapter)
				rq.Equal(tc.expHighlight.Location, highlight.Location)
				rq.Equal(tc.expHighlight.URL, highlight.URL)
			}
		})
	}
}

func TestSqlite_DeleteHighlight(t *testing.T) {
	tt := []struct {
		name   string
		id     int
		setup  func(*testing.T, sqlmock.Sqlmock)
		expErr bool
	}{
		{
			name: "success",
			id:   1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`
					SELECT book_id 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(2))

				mock.ExpectExec(`
					DELETE FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec(`
					UPDATE books 
					SET updated = ?
					WHERE id = ?
				`).WithArgs(sqlmock.AnyArg(), 2).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			expErr: false,
		},
		{
			name: "query error",
			id:   1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`
					SELECT book_id 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expErr: true,
		},
		{
			name: "scan error",
			id:   1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`
					SELECT book_id 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow("invalid"))
				mock.ExpectRollback()
			},
			expErr: true,
		},
		{
			name: "delete error",
			id:   1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`
					SELECT book_id 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(2))

				mock.ExpectExec(`
					DELETE FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expErr: true,
		},
		{
			name: "update error",
			id:   1,
			setup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`
					SELECT book_id 
					FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(2))

				mock.ExpectExec(`
					DELETE FROM highlights 
					WHERE id = ?
				`).WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec(`
					UPDATE books 
					SET updated = ?
					WHERE id = ?
				`).WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			db, mock, err := sqlmock.New(
				sqlmock.QueryMatcherOption(queryMatcher(t)),
			)
			rq.NoError(err)

			s := &DB{
				db: db,
			}

			tc.setup(t, mock)

			if tc.expErr {
				err := s.DeleteHighlight(context.Background(), tc.id)
				rq.Error(err)
			} else {
				err := s.DeleteHighlight(context.Background(), tc.id)
				rq.NoError(err)
			}
		})
	}
}
