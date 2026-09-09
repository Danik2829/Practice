package storage

import (
	"database/sql"
	"fmt"
	"restAPI/config"
	"restAPI/core"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewDbStorage(cfg config.Config) (*Store, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Create(u core.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (id, first_name, last_name, email, age) VALUES ($1, $2, $3, $4, $5)",
		u.ID, u.FirstName, u.LastName, u.Email, u.Age,
	)
	return err
}

func (s *Store) Get(id string) *core.User {
	var user core.User
	err := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return nil
	} else {
		return &user
	}
}

func (s *Store) Update(user core.User) (core.User, error) {
	var updated core.User
	err := s.db.QueryRow(`
        UPDATE users 
        SET first_name = $2, last_name = $3, email = $4, age = $5 
        WHERE id = $1
        RETURNING id, first_name, last_name, email, age
    `, user.ID, user.FirstName, user.LastName, user.Email, user.Age).
		Scan(&updated.ID, &updated.FirstName, &updated.LastName, &updated.Email, &updated.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return core.User{}, core.NotFound
		}
		return core.User{}, err
	}
	return updated, nil
}

func (s *Store) Delete(id string) string {
	var deletedId string
	err := s.db.QueryRow("DELETE FROM users WHERE id = $1 RETURNING id", id).Scan(&deletedId)
	if err != nil {
		return ""
	}
	if deletedId == "" {
		return ""
	} else {
		return deletedId
	}
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
