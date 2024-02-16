package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"wishlist_auth/internal/http-server/errors"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, password, dbName string) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{db: db}

	cwd, _ := os.Getwd()
	log.Println("Current working directory:", cwd)

	err = goose.Up(storage.db, "db/migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
}

func (s *Storage) CreateUser(lname, fname, email string, passwordHashed []byte) (int, error) {
	const op = "storage.postgres.CreateUser"

	query := `
		INSERT INTO wish_users (lname, fname, email, pwd_hash) VALUES ($1, $2, $3, $4) RETURNING uid;
		`

	var id int

	err := s.db.QueryRow(query, lname, fname, email, passwordHashed).Scan(&id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return 0, errors.ErrUsernameTaken
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Login(email string) ([]byte, int, error) {
	const op = "storage.postgres.Login"

	query := `
		SELECT wish_users.pwd_hash, wish_users.uid FROM wish_users WHERE email = $1 LIMIT 1;
		`

	var hash []byte
	var id int

	err := s.db.QueryRow(query, email).Scan(&hash, &id)
	if err == sql.ErrNoRows {
		return nil, 0, errors.ErrIncorrectEmail
	} else if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, 0, errors.ErrUsernameTaken
		}
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return hash, id, nil
}
