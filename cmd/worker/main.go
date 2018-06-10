package main

import (
	"flag"
	"log"
	"os"

	"gitlab.com/luizbranco/srs/db/psql"
	"gitlab.com/luizbranco/srs/worker"
)

var dbURL string

func init() {
	dbURL = os.Getenv("DATABASE_URL")

	if dbURL == "" {
		dbURL = "postgres://srs:s3cr3t@192.168.0.11:5432/srs?sslmode=disable"
	}

	flag.Parse()
}

func main() {
	db, err := psql.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool := &worker.Worker{
		Database: db,
	}

	log.Fatal(pool.Start())
}
