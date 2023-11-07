package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
	"url-shrotener/internal/config"
)

type Database struct {
	DB *sql.DB
}

func NewDB(cfg *config.Config) (*Database, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.DBPassword,
		cfg.DB.DBName,
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	time.Sleep(5 * time.Second)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (s *Database) SaveURL(urlToSave string, alias string) error {
	const op = "storage.db.SaveURL"
	stmt, err := s.DB.Prepare("INSERT INTO urls(full_url, short_url) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "urls_full_url_key" {
				return fmt.Errorf("%s: %w", op, ErrURLExists)
			}

			if pqErr.Constraint == "urls_short_url_key" {
				return fmt.Errorf("%s: %w", op, AliasExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Database) GetURL(alias string) (string, error) {
	const op = "storage.db.GetURL"
	stmt, err := s.DB.Prepare("SELECT full_url FROM urls WHERE short_url = $1 LIMIT 1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return url, nil
}

func (s *Database) GetAlias(url string) (string, error) {
	const op = "storage.db.GetAlias"

	stmt, err := s.DB.Prepare("SELECT short_url FROM urls WHERE full_url = $1 LIMIT 1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var alias string
	err = stmt.QueryRow(url).Scan(&alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return alias, nil
}
