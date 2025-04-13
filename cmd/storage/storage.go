package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	conf "github.com/LysetsDal/webscraper-v2/config"
	T "github.com/LysetsDal/webscraper-v2/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateEntry(*T.ListMessage) error
	DeleteEntry(int) error
}

type PostgresStore struct {
	db *sql.DB
}

type DuplicateKeyError struct {
}

func NewPostgresStore() (*PostgresStore, error) {
	godotenv.Load(conf.ENV_FILE_PATH) // Gets secret from .env
	connStr := os.Getenv(conf.POSTGRES_DB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}

func (s PostgresStore) Init() error {
	return s.CreateOpenListTable()
}

func (s *PostgresStore) CreateOpenListTable() error {
	query := `CREATE TABLE IF NOT EXISTS openlist (
		id SERIAL PRIMARY KEY NOT NULL,
		name VARCHAR(100) UNIQUE,
		state VARCHAR(10),
		updated_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) InsertEntry(e T.Entry) error {
	query := `INSERT INTO openlist 
	(name, state, updated_at) 
	VALUES ($1, $2, $3)`

	_, err := s.db.Exec(query, e.Name, e.Status, e.UpdatedAt)
	return err
}

func (s *PostgresStore) InsertEntries(es []T.Entry) error {
	query := `INSERT INTO openlist 
	(name, state, updated_at) 
	VALUES `

	values := []interface{}{}
	placeholders := []string{}

	for i, e := range es {
		idx := i * 3
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d)", idx+1, idx+2, idx+3),
		)
		values = append(values, e.Name, e.Status, e.UpdatedAt)
	}

	query += strings.Join(placeholders, ", ")
	if _, err := s.db.Exec(query, values...); err != nil {
		return fmt.Errorf("InsertEntries failed db insert query")
	}

	return nil
}
