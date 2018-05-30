package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luizbranco/srs/web/html"
	"github.com/luizbranco/srs/web/psql"
	"github.com/luizbranco/srs/web/server"
)

func main() {
	db, err := psql.New("192.168.0.11", "5432", "srs", "srs", "s3cr3t")
	if err != nil {
		log.Fatal(err)
	}

	srv := &server.Server{
		Template: html.New("web/templates"),
		Database: db,
	}
	mux := srv.NewServeMux()

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
