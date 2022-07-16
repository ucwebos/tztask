package conf

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DBFile string `json:"db_file" yaml:"db_file"`
}

// CreateInstance .
func (s *SQLite) CreateInstance() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", s.DBFile)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}
