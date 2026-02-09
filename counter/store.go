package counter

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("counter db path is empty")
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create counter db dir: %w", err)
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open counter db: %w", err)
	}

	store := &Store{db: db}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) init() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("counter store is nil")
	}

	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS visitors (
	ip TEXT PRIMARY KEY,
	first_seen INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS opt_out (
	ip TEXT PRIMARY KEY,
	opted_out_at INTEGER NOT NULL
);
`)
	if err != nil {
		return fmt.Errorf("init counter db: %w", err)
	}
	return nil
}

func (s *Store) IsOptedOut(ip string) (bool, error) {
	if s == nil || s.db == nil {
		return false, fmt.Errorf("counter store is nil")
	}
	if ip == "" {
		return false, nil
	}

	var exists int
	if err := s.db.QueryRow(`SELECT 1 FROM opt_out WHERE ip = ? LIMIT 1;`, ip).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("read opt-out: %w", err)
	}
	return true, nil
}

func (s *Store) RecordVisit(ip string) (int, error) {
	if s == nil || s.db == nil {
		return 0, fmt.Errorf("counter store is nil")
	}

	optedOut, err := s.IsOptedOut(ip)
	if err != nil {
		return 0, err
	}
	if optedOut {
		return s.Count()
	}

	if ip != "" {
		if _, err := s.db.Exec(`
INSERT OR IGNORE INTO visitors (ip, first_seen)
VALUES (?, strftime('%s','now'));
`, ip); err != nil {
			return 0, fmt.Errorf("record visit: %w", err)
		}
	}

	return s.Count()
}

func (s *Store) SetOptOut(ip string, optOut bool) (int, error) {
	if s == nil || s.db == nil {
		return 0, fmt.Errorf("counter store is nil")
	}
	if ip == "" {
		return s.Count()
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin privacy tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if optOut {
		if _, err := tx.Exec(`
INSERT OR IGNORE INTO opt_out (ip, opted_out_at)
VALUES (?, strftime('%s','now'));
`, ip); err != nil {
			return 0, fmt.Errorf("opt-out insert: %w", err)
		}
		if _, err := tx.Exec(`DELETE FROM visitors WHERE ip = ?;`, ip); err != nil {
			return 0, fmt.Errorf("opt-out delete: %w", err)
		}
	} else {
		if _, err := tx.Exec(`DELETE FROM opt_out WHERE ip = ?;`, ip); err != nil {
			return 0, fmt.Errorf("opt-out clear: %w", err)
		}
		if _, err := tx.Exec(`
INSERT OR IGNORE INTO visitors (ip, first_seen)
VALUES (?, strftime('%s','now'));
`, ip); err != nil {
			return 0, fmt.Errorf("opt-in insert: %w", err)
		}
	}

	var count int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM visitors;`).Scan(&count); err != nil {
		return 0, fmt.Errorf("read counter: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit privacy tx: %w", err)
	}

	return count, nil
}

func (s *Store) Count() (int, error) {
	if s == nil || s.db == nil {
		return 0, fmt.Errorf("counter store is nil")
	}

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM visitors;`).Scan(&count); err != nil {
		return 0, fmt.Errorf("read counter: %w", err)
	}
	return count, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
