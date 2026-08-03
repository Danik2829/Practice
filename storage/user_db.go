package storage

import (
	"database/sql"
	"fmt"
	"os"
	"restAPI/core"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewDbStorage() (*Store, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("No .env file found; falling back to system environment variables")
	}
	connStr := fmt.Sprintf(`user=%s password=%s dbname=%s sslmode=%s`,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))
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

func (s *Store) Get(id string) (core.User, error) {
	var user core.User
	err := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return core.User{}, err
	} else {
		return user, nil
	}
}

func (s *Store) Update(user core.User) error {
	result, err := s.db.Exec(`UPDATE users SET 
	firstName = $2,
	lastName = $3, 
	email = $4, 
	age = $5 
	WHERE id = $1`,
		user.ID, user.FirstName, user.LastName, user.Email, user.Age)
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

func (s *Store) Delete(id string) error {
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

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
