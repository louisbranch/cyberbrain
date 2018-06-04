package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luizbranco/srs/db/psql"
	"github.com/luizbranco/srs/generator"
	"github.com/luizbranco/srs/web/html"
	"github.com/luizbranco/srs/web/server"
	"github.com/luizbranco/srs/web/urlbuilder"
)

func main() {
	db, err := psql.New("192.168.0.11", "5432", "srs", "srs", "s3cr3t")
	if err != nil {
		log.Fatal(err)
	}

	ub, err := urlbuilder.New()
	if err != nil {
		log.Fatal(err)
	}

	gen := generator.Generator{}
	tpl := html.New("web/templates")

	srv := &server.Server{
		Template:          tpl,
		Database:          db,
		URLBuilder:        ub,
		PracticeGenerator: gen,
	}

	mux := srv.NewServeMux()

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
