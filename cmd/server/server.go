package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luizbranco/srs/web/html"
	"github.com/luizbranco/srs/web/server"
	"github.com/luizbranco/srs/web/sqlite"
)

func main() {
	db, err := sqlite.New("srs.db")

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
