package main

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/luizbranco/srs/authentication"
	"gitlab.com/luizbranco/srs/db/psql"
	"gitlab.com/luizbranco/srs/generator"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server"
	"gitlab.com/luizbranco/srs/web/session"
	"gitlab.com/luizbranco/srs/web/urlbuilder"
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

	gen := generator.Generator{
		Database: db,
	}

	auth := authentication.Authenticator{}

	session := &session.Manager{
		Database: db,
		Secret:   "saahskdhsakjdao8oKAAKJJAJSkjasEE", // FIXME
	}

	tpl := html.New("web/templates")

	srv := &server.Server{
		Template:          tpl,
		Database:          db,
		URLBuilder:        ub,
		PracticeGenerator: gen,
		Authenticator:     auth,
		SessionManager:    session,
	}

	mux := srv.NewServeMux()

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
