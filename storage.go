package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreatePlace(*Place) error
	DeletePlace(int) error
	UpdatePlace(*Place) error
	GetPlaceByID(int) (*Place, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// TODO set up envvars for database connection
	connStr := "user=postgres dbname=postgres password=goplaces sslmode=disable"
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

// TODO build out PostgresStore methods
func (s *PostgresStore) CreatePlace(*Place) error {
	return nil
}

func (s *PostgresStore) DeletePlace(id int) error {
	return nil
}

func (s *PostgresStore) UpdatePlace(*Place) error {
	return nil
}

func (s *PostgresStore) GetPlaceByID(id int) (*Place, error) {
	return nil, nil
}
