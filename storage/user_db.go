package storage

import (
	"database/sql"
	"fmt"
	"restAPI/core"

	_ "github.com/lib/pq"
)

type DbStorage struct {
	db *sql.DB
}

func NewDbStorage() (*DbStorage, error) {
	connStr := "user=postgres password=21282908 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	return &DbStorage{db: db}, nil
}

func (s *DbStorage) Create(u core.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (id, first_name, last_name, email, age) VALUES ($1, $2, $3, $4, $5)",
		u.ID, u.FirstName, u.LastName, u.Email, u.Age,
	)
	return err
}

func (s *DbStorage) Get(id string) (core.User, error) {
	var user core.User
	err := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return core.User{}, err
	} else {
		return user, nil
	}
}

func (s *DbStorage) Update(user core.User) error {
	result, err := s.db.Exec("UPDATE users SET firstName = $2, lastName = $3, email = $4, age = $5 WHERE id = $1", user.ID, user.FirstName, user.LastName, user.Email, user.Age)
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return core.NotFound
	}
	return nil
}

func (s *DbStorage) Delete(id string) error {
	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return core.NotFound
	}
	return nil
}

func (s *DbStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
