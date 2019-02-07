package db_client

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func (d *DB) Start() {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	d.db = db
}
