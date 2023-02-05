package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreatePlace(*Place) error
	DeletePlace(int) error
	UpdatePlace(*Place) error
	GetPlaces() ([]*Place, error)
	GetPlaceByID(int) (*Place, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("An error occured. Err: %s", err)
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
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

func (s *PostgresStore) Init() error {
	return s.createPlaceTable()
}

func (s *PostgresStore) createPlaceTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS place (
		id serial primary key,
		location varchar(50),
		name varchar(50),
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreatePlace(place *Place) error {
	sqlStatement := `
	INSERT INTO place (location, name, created_at)
	VALUES ($1, $2, $3)`

	resp, err := s.db.Exec(sqlStatement, place.Location, place.Name, place.CreatedAt)

	fmt.Println(resp)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeletePlace(id int) error {
	// TODO In production it would be better to soft delete - flag as deleted and make unavailable via the API
	_, err := s.db.Query("DELETE FROM place WHERE id = $1", id)
	return err
}

// TODO build UpdatePlace function
func (s *PostgresStore) UpdatePlace(*Place) error {
	return nil
}

func (s *PostgresStore) GetPlaceByID(id int) (*Place, error) {
	rows, err := s.db.Query("SELECT * FROM place WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoPlace(rows)
	}
	return nil, fmt.Errorf("Place %d not found", id)
}

func (s *PostgresStore) GetPlaces() ([]*Place, error) {
	rows, err := s.db.Query("SELECT * FROM place")
	if err != nil {
		return nil, err
	}

	places := []*Place{}
	for rows.Next() {
		place, err := scanIntoPlace(rows)
		if err != nil {
			return nil, err
		}
		places = append(places, place)
	}

	return places, nil
}

func scanIntoPlace(rows *sql.Rows) (*Place, error) {
	place := new(Place)
	err := rows.Scan(&place.ID, &place.Location, &place.Name, &place.CreatedAt)

	return place, err
}
