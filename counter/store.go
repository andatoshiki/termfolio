package counter

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/oschwald/geoip2-golang"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB

	statsMu     sync.Mutex
	statsDirty  bool
	statsPath   string
	statsCache  CountryStats
	statsReader *geoip2.Reader
}

type CountryStats struct {
	TotalVisitors      int
	TopCountry         string
	TopCountryVisitors int
	TopCountries       []CountryCount
}

type CountryCount struct {
	Name     string
	Visitors int
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

	store := &Store{
		db:         db,
		statsDirty: true,
	}
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
		result, err := s.db.Exec(`
INSERT OR IGNORE INTO visitors (ip, first_seen)
VALUES (?, strftime('%s','now'));
`, ip)
		if err != nil {
			return 0, fmt.Errorf("record visit: %w", err)
		}
		if rowsChanged(result) {
			s.invalidateStatsCache()
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
	visitorChanged := false

	if optOut {
		if _, err := tx.Exec(`
INSERT OR IGNORE INTO opt_out (ip, opted_out_at)
VALUES (?, strftime('%s','now'));
`, ip); err != nil {
			return 0, fmt.Errorf("opt-out insert: %w", err)
		}
		delResult, err := tx.Exec(`DELETE FROM visitors WHERE ip = ?;`, ip)
		if err != nil {
			return 0, fmt.Errorf("opt-out delete: %w", err)
		}
		visitorChanged = rowsChanged(delResult)
	} else {
		if _, err := tx.Exec(`DELETE FROM opt_out WHERE ip = ?;`, ip); err != nil {
			return 0, fmt.Errorf("opt-out clear: %w", err)
		}
		insResult, err := tx.Exec(`
INSERT OR IGNORE INTO visitors (ip, first_seen)
VALUES (?, strftime('%s','now'));
`, ip)
		if err != nil {
			return 0, fmt.Errorf("opt-in insert: %w", err)
		}
		visitorChanged = rowsChanged(insResult)
	}

	var count int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM visitors;`).Scan(&count); err != nil {
		return 0, fmt.Errorf("read counter: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit privacy tx: %w", err)
	}

	if visitorChanged {
		s.invalidateStatsCache()
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

func (s *Store) CountryStats(geoLiteDBPath string) (CountryStats, error) {
	if s == nil || s.db == nil {
		return CountryStats{}, fmt.Errorf("counter store is nil")
	}
	if strings.TrimSpace(geoLiteDBPath) == "" {
		return CountryStats{}, fmt.Errorf("geolite db path is empty")
	}
	if _, err := os.Stat(geoLiteDBPath); err != nil {
		if os.IsNotExist(err) {
			return CountryStats{}, fmt.Errorf("geolite db not found")
		}
		return CountryStats{}, fmt.Errorf("geolite db is not accessible")
	}

	s.statsMu.Lock()
	defer s.statsMu.Unlock()

	if geoLiteDBPath != s.statsPath {
		if s.statsReader != nil {
			_ = s.statsReader.Close()
			s.statsReader = nil
		}
		s.statsPath = geoLiteDBPath
		s.statsDirty = true
	}

	if !s.statsDirty {
		return s.statsCache, nil
	}

	if s.statsReader == nil {
		reader, err := geoip2.Open(geoLiteDBPath)
		if err != nil {
			return CountryStats{}, fmt.Errorf("geolite db is invalid or unreadable")
		}
		s.statsReader = reader
	}

	totalVisitors, err := s.Count()
	if err != nil {
		return CountryStats{}, err
	}

	rows, err := s.db.Query(`SELECT ip FROM visitors;`)
	if err != nil {
		return CountryStats{}, fmt.Errorf("query visitor IPs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	countryCounts := make(map[string]int)
	for rows.Next() {
		var ipValue string
		if err := rows.Scan(&ipValue); err != nil {
			return CountryStats{}, fmt.Errorf("scan visitor IP: %w", err)
		}

		parsedIP := net.ParseIP(strings.TrimSpace(ipValue))
		if parsedIP == nil {
			continue
		}

		record, err := s.statsReader.Country(parsedIP)
		if err != nil || record == nil {
			continue
		}

		country := strings.TrimSpace(record.Country.Names["en"])
		if country == "" {
			country = strings.TrimSpace(record.Country.IsoCode)
		}
		if country == "" {
			continue
		}
		countryCounts[country]++
	}

	if err := rows.Err(); err != nil {
		return CountryStats{}, fmt.Errorf("iterate visitor IPs: %w", err)
	}

	topCountries := make([]CountryCount, 0, len(countryCounts))
	for country, count := range countryCounts {
		topCountries = append(topCountries, CountryCount{
			Name:     country,
			Visitors: count,
		})
	}

	sort.Slice(topCountries, func(i, j int) bool {
		if topCountries[i].Visitors == topCountries[j].Visitors {
			return topCountries[i].Name < topCountries[j].Name
		}
		return topCountries[i].Visitors > topCountries[j].Visitors
	})

	if len(topCountries) > 5 {
		topCountries = topCountries[:5]
	}

	topCountry := "N/A"
	topCountryVisitors := 0
	if len(topCountries) > 0 {
		topCountry = topCountries[0].Name
		topCountryVisitors = topCountries[0].Visitors
	}

	stats := CountryStats{
		TotalVisitors:      totalVisitors,
		TopCountry:         topCountry,
		TopCountryVisitors: topCountryVisitors,
		TopCountries:       append([]CountryCount(nil), topCountries...),
	}
	s.statsCache = stats
	s.statsDirty = false
	return stats, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	s.statsMu.Lock()
	if s.statsReader != nil {
		_ = s.statsReader.Close()
		s.statsReader = nil
	}
	s.statsMu.Unlock()
	return s.db.Close()
}

func (s *Store) invalidateStatsCache() {
	if s == nil {
		return
	}
	s.statsMu.Lock()
	s.statsDirty = true
	s.statsMu.Unlock()
}

func rowsChanged(result sql.Result) bool {
	if result == nil {
		return false
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return true
	}
	return changed > 0
}
