package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type database struct {
	db *sql.DB
}

func (s *database) connect() {
	connStr := "user=secretanry password=2271799 host=158.160.73.176 port=128 database=gitdiff sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
	res, err := s.db.Exec("SELECT current_user")
	if err != nil {
		fmt.Print(res)
	}
}

func (s *database) connector() {
	if s.db == nil {
		s.connect()
	}
}

func (s *database) createDifferenceTable() {
	_, err := s.db.Exec("select * from diff")
	if err != nil {
		_, err = s.db.Exec("create table diff (hash1 varchar(255) not null, " +
			"hash2 varchar(255) not null, difference jsonb not null)")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *database) createInfoTable() {
	_, err := s.db.Exec("select * from info")
	if err != nil {
		_, err := s.db.Exec("create table info (hash1 varchar(255) not null, " +
			"message varchar(255), branch varchar(255), primary key(hash1))")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *database) dropTables() {
	_, err := s.db.Exec("select * from diff")
	if err == nil {
		_, err = s.db.Exec("drop table diff")
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = s.db.Exec("select * from info")
	if err == nil {
		_, err = s.db.Exec("drop table info")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *database) addDiffData(arr []string) {
	s.createDifferenceTable()
	_, err := s.db.Prepare("INSERT into diff values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = s.db.Exec(arr[0], arr[1], arr[2])
	if err != nil {
		log.Fatal(err)
	}
}

func (s *database) addCommitInfo(arr []string) {
	s.createInfoTable()
	res, err := s.db.Query("SELECT * where hash=" + arr[0])
	if err != nil {
		log.Fatal(err)
	}
	if !res.Next() {
		_, err := s.db.Prepare("INSERT into info values (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = s.db.Exec(arr[0], arr[1], arr[2])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *database) getDiff(hash1 string, hash2 string) string {
	s.createDifferenceTable()
	res, err := s.db.Query("select difference where hash1=" + hash1 + " and hash2=" + hash2)
	if err != nil {
		log.Fatal(err)
	}
	r, err := res.Columns()
	return r[0]
}
