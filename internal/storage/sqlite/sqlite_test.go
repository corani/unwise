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
		name    string
		errPing error
		err1    error
		err2    error
		expErr  bool
	}{
		{
			name:    "success",
			errPing: nil,
			err1:    nil,
			err2:    nil,
			expErr:  false,
		},
		{
			name:    "ping error",
			errPing: assert.AnError,
			err1:    nil,
			err2:    nil,
			expErr:  true,
		},
		{
			name:    "error 1",
			errPing: nil,
			err1:    assert.AnError,
			err2:    nil,
			expErr:  true,
		},
		{
			name:    "error 2",
			errPing: nil,
			err1:    nil,
			err2:    assert.AnError,
			expErr:  true,
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

			mock.ExpectPing().
				WillReturnError(tc.errPing)

			exec := mock.ExpectExec(`
				DROP TABLE IF EXISTS books; 
				DROP TABLE IF EXISTS highlights;
			`)

			if tc.err1 != nil {
				exec.WillReturnError(tc.err1)
			} else {
				exec.WillReturnResult(sqlmock.NewResult(0, 0))
			}

			exec = mock.ExpectExec(`
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
			`)

			if tc.err2 != nil {
				exec.WillReturnError(tc.err2)
			} else {
				exec.WillReturnResult(sqlmock.NewResult(0, 0))
			}

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
		name     string
		title    string
		author   string
		source   string
		err      error
		rows     *sqlmock.Rows
		expTitle string
		expErr   bool
		expBook  storage.Book
	}{
		{
			name:   "success",
			title:  "title",
			author: "author",
			source: "source",
			err:    nil,
			rows: sqlmock.NewRows([]string{"id", "count"}).
				AddRow(1, 2),
			expTitle: "title",
			expErr:   false,
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
			err:    nil,
			rows: sqlmock.NewRows([]string{"id", "count"}).
				AddRow(1, 2),
			expTitle: storage.DefaultTitle,
			expErr:   false,
			expBook: storage.Book{
				ID:            1,
				Title:         storage.DefaultTitle,
				Author:        "author",
				SourceURL:     "source",
				NumHighlights: 2,
			},
		},
		{
			name:     "query error",
			title:    "title",
			author:   "author",
			source:   "source",
			err:      assert.AnError,
			expTitle: "title",
			expErr:   true,
		},
		{
			name:   "scan error",
			title:  "title",
			author: "author",
			source: "source",
			err:    nil,
			rows: sqlmock.NewRows([]string{"id", "count"}).
				AddRow("invalid", "invalid"),
			expTitle: "title",
			expErr:   true,
		},
		{
			name:     "no rows",
			title:    "title",
			author:   "author",
			source:   "source",
			err:      nil,
			rows:     sqlmock.NewRows([]string{"id", "count"}),
			expTitle: "title",
			expErr:   true,
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

			exec := mock.ExpectQuery(`
				INSERT INTO books (title, author, source_url, updated) 
				VALUES (?, ?, ?, ?) 
				ON CONFLICT (title, author, source_url) DO UPDATE SET updated = ?
				RETURNING id, (SELECT COUNT(*) FROM highlights WHERE book_id = id)
			`).WithArgs(tc.expTitle, tc.author, tc.source, sqlmock.AnyArg(), sqlmock.AnyArg())

			if tc.err != nil {
				exec.WillReturnError(tc.err)
			} else {
				exec.WillReturnRows(tc.rows)
			}

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
		err      error
		rows     *sqlmock.Rows
		expErr   bool
		expBooks []storage.Book
	}{
		{
			name: "success",
			lt:   now,
			gt:   now,
			err:  nil,
			rows: sqlmock.NewRows([]string{"id", "title", "author", "source_url", "updated", "num_highlights"}).
				AddRow(1, "title", "author", "source", now, 2),
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
			name:   "query error",
			lt:     now,
			gt:     now,
			err:    assert.AnError,
			rows:   nil,
			expErr: true,
		},
		{
			name: "scan error",
			lt:   now,
			gt:   now,
			err:  nil,
			rows: sqlmock.NewRows([]string{"id", "title", "author", "source_url", "updated", "num_highlights"}).
				AddRow("invalid", "title", "author", "source", now, 2),
			expErr: true,
		},
		{
			name: "time error",
			lt:   now,
			gt:   now,
			err:  nil,
			rows: sqlmock.NewRows([]string{"id", "title", "author", "source_url", "updated", "num_highlights"}).
				AddRow(1, "title", "author", "source", "invalid", 2),
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

			exec := mock.ExpectQuery(`
				SELECT b.id, b.title, b.author, b.source_url, b.updated, COUNT(h.id) AS num_highlights
				FROM   books AS b LEFT OUTER JOIN highlights AS h ON b.id = h.book_id
				WHERE  b.updated >= ? AND b.updated <= ?
				GROUP BY b.id
			`).WithArgs(tc.lt.Format(time.RFC3339), tc.gt.Format(time.RFC3339))

			if tc.err != nil {
				exec.WillReturnError(tc.err)
			} else {
				exec.WillReturnRows(tc.rows)
			}

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
		err          error
		rows         *sqlmock.Rows
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
			err:      nil,
			rows:     sqlmock.NewRows([]string{"id"}).AddRow(1),
			expErr:   false,
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
			err:      assert.AnError,
			rows:     nil,
			expErr:   true,
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
			err:      nil,
			rows:     sqlmock.NewRows([]string{"id"}).AddRow("invalid"),
			expErr:   true,
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
			err:      nil,
			rows:     sqlmock.NewRows([]string{"id"}),
			expErr:   true,
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

			exec := mock.ExpectQuery(`
				INSERT INTO highlights (book_id, text, note, chapter, location, url, updated)
				VALUES (?, ?, ?, ?, ?, ?, ?) 
				ON CONFLICT (book_id, text) DO UPDATE SET note = ?, chapter = ?, location = ?, url = ?, updated = ?
				RETURNING id
			`).WithArgs(
				tc.book.ID, tc.text, tc.note, tc.chapter, tc.location, tc.url, now.Format(time.RFC3339),
				tc.note, tc.chapter, tc.location, tc.url, now.Format(time.RFC3339),
			)

			if tc.err != nil {
				exec.WillReturnError(tc.err)
			} else {
				exec.WillReturnRows(tc.rows)
			}

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
		err           error
		rows          *sqlmock.Rows
		expErr        bool
		expHighlights []storage.Highlight
	}{
		{
			name:   "success",
			bookID: 1,
			err:    nil,
			rows: sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
				AddRow(1, 1, "text1", "note1", "chapter1", 10, "url1", now).
				AddRow(2, 1, "text2", "note2", "chapter2", 20, "url2", now),
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
			err:    assert.AnError,
			rows:   nil,
			expErr: true,
		},
		{
			name:   "scan error",
			bookID: 1,
			err:    nil,
			rows: sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
				AddRow("invalid", 1, "text", "note", "chapter", 3, "url", now),
			expErr: true,
		},
		{
			name:   "time error",
			bookID: 1,
			err:    nil,
			rows: sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
				AddRow(1, 1, "text", "note", "chapter", 3, "url", "invalid"),
			expErr: true,
		},
		{
			name:          "no highlights",
			bookID:        999,
			err:           nil,
			rows:          sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}),
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

			exec := mock.ExpectQuery(`
				SELECT id, book_id, text, note, chapter, location, url, updated 
				FROM highlights 
				WHERE book_id = ? 
				ORDER BY location
			`).WithArgs(tc.bookID)

			if tc.err != nil {
				exec.WillReturnError(tc.err)
			} else {
				exec.WillReturnRows(tc.rows)
			}

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
		err           error
		rows          *sqlmock.Rows
		expErr        bool
		expHighlights []storage.Highlight
	}{
		{
			name: "success",
			lt:   now,
			gt:   now,
			err:  nil,
			rows: sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
				AddRow(1, 1, "text", "note", "chapter", 3, "url", now),
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
			name:   "query error",
			lt:     now,
			gt:     now,
			err:    assert.AnError,
			rows:   nil,
			expErr: true,
		},
		{
			name: "scan error",
			lt:   now,
			gt:   now,
			err:  nil,
			rows: sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
				AddRow("invalid", 1, "text", "note", "chapter", 3, "url", now),
			expErr: true,
		},
		{
			name: "time error",
			lt:   now,
			gt:   now,
			err:  nil,
			rows: sqlmock.NewRows([]string{"id", "book_id", "text", "note", "chapter", "location", "url", "updated"}).
				AddRow(1, 1, "text", "note", "chapter", 3, "url", "invalid"),
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

			exec := mock.ExpectQuery(`
				SELECT id, book_id, text, note, chapter, location, url, updated 
				FROM   highlights 
				WHERE  updated >= ? AND updated <= ? 
			`).WithArgs(tc.lt.Format(time.RFC3339), tc.gt.Format(time.RFC3339))

			if tc.err != nil {
				exec.WillReturnError(tc.err)
			} else {
				exec.WillReturnRows(tc.rows)
			}

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
