package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type database struct {
	db        *sql.DB
	diffTable bool
	infoTable bool
}

func (s *database) connect() {
	connStr := "user=secretanry password=2271799 host=localhost port=5432 database=gitdiff sslmode=disable"
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

func (s *database) createDifferenceTable() {
	if !s.diffTable {
		_, err := s.db.Exec("create table diff (hash1 varchar(255) not null, " +
			"hash2 varchar(255) not null, difference jsonb not null)")
		if err != nil {
			log.Fatal(err)
		}
		s.diffTable = true
	}
}

func (s *database) createInfoTable() {
	if !s.infoTable {
		_, err := s.db.Exec("create table info (hash1 varchar(255) not null, " +
			"message varchar(255), primary key(hash1))")
		if err != nil {
			log.Fatal(err)
		}
		s.infoTable = true
	}
}

func (s *database) dropTables() {
	_, err := s.db.Exec("drop table diff")
	if err != nil {
		log.Fatal(err)
	}
	s.diffTable = false
	_, err = s.db.Exec("drop table info")
	if err != nil {
		log.Fatal(err)
	}
	s.infoTable = false
}
