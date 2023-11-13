package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"tgbot-read-adviser/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	s := Storage{db: db}
	if err := s.init(context.TODO()); err != nil {
		return nil, fmt.Errorf("can't init database: %w", err)
	}

	return &s, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages(url, user_name) VALUES($1,$2)`
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = $1 ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	isExist, err := s.IsExists(ctx, page)
	if err != nil {
		fmt.Errorf("can't check exists file: %w", err)
	}

	if !isExist {
		return storage.ErrNoSavedPages
	}

	q := `DELETE FROM pages WHERE url = $1 AND user_name = $2`
	_, err = s.db.ExecContext(ctx, q, page.URL, page.UserName)
	if err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (s *Storage) All(ctx context.Context, userName string) ([]*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = $1`
	rows, err := s.db.QueryContext(ctx, q, userName)
	if err != nil {
		return nil, fmt.Errorf("can't get all pages: %w", err)
	}

	result := make([]*storage.Page, 0, 4)
	var url string
	i := 0

	for rows.Next() {
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("cna't scan data from rows: %w", err)
		}
		result = append(result, &storage.Page{URL: url, UserName: userName})
		i++
	}

	if i == 0 {
		return nil, storage.ErrNoSavedPages
	}

	return result, nil
}

func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = $1 AND user_name = $2`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages(url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
