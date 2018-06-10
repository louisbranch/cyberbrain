package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/luizbranco/srs/authentication"
	"gitlab.com/luizbranco/srs/db/psql"
	"gitlab.com/luizbranco/srs/generator"
	"gitlab.com/luizbranco/srs/web/html"
	"gitlab.com/luizbranco/srs/web/server"
	"gitlab.com/luizbranco/srs/web/session"
	"gitlab.com/luizbranco/srs/web/urlbuilder"
)

var dbURL, sessionSecret, hashidSalt string

func init() {
	dbURL = os.Getenv("DATABASE_URL")
	sessionSecret = os.Getenv("SESSION_SECRET")
	hashidSalt = os.Getenv("HASHID_SALT")

	if dbURL == "" {
		dbURL = "postgres://srs:s3cr3t@192.168.0.11:5432/srs?sslmode=disable"
	}

	if sessionSecret == "" {
		sessionSecret = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	}

	if hashidSalt == "" {
		hashidSalt = "s3cret"
	}

	flag.StringVar(&dbURL, "database-url", dbURL, "database connection url")
	flag.StringVar(&sessionSecret, "session-secret", sessionSecret, "session cookie id secret")
	flag.StringVar(&hashidSalt, "hashid-salt", hashidSalt, "salt for hashid url")

	flag.Parse()
}

func main() {
	db, err := psql.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	ub, err := urlbuilder.New(hashidSalt)
	if err != nil {
		log.Fatal(err)
	}

	gen := generator.Generator{
		Database: db,
	}

	auth := authentication.Authenticator{}

	session := &session.Manager{
		Database: db,
		Secret:   sessionSecret,
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
