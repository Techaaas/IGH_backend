package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type database struct {
	db *sql.DB
}

func (s *database) connect() {
	connStr := "user=secretanry password=2271799 host=jdbc:postgresql://localhost:5432/postgres sslmode=disabled"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
}

func (s *database) connector() {
	if s.db == nil {
		s.connect()
	}
}
