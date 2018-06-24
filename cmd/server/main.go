package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/luizbranco/cyberbrain/authentication"
	"gitlab.com/luizbranco/cyberbrain/db/psql"
	"gitlab.com/luizbranco/cyberbrain/generator"
	"gitlab.com/luizbranco/cyberbrain/web/html"
	"gitlab.com/luizbranco/cyberbrain/web/server"
	"gitlab.com/luizbranco/cyberbrain/web/session"
	"gitlab.com/luizbranco/cyberbrain/web/urlbuilder"
	"gitlab.com/luizbranco/cyberbrain/worker"
	"gitlab.com/luizbranco/cyberbrain/worker/resizer"
)

const (
	Development = "dev"
	Production  = "prod"
)

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	dbURL := os.Getenv("DATABASE_URL")
	sessionSecret := os.Getenv("SESSION_SECRET")
	hashidSalt := os.Getenv("HASHID_SALT")
	env := os.Getenv("ENVIRONMENT")

	if env == "" {
		env = Development
	}

	blitlineID := os.Getenv("BLITLINE_ID")
	blitlineCallbackURL := os.Getenv("BLITLINE_CALLBACK_URL")

	awsBucket := os.Getenv("AWS_BUCKET")

	piioDomain := os.Getenv("PIIO_DOMAIN")
	piioID := os.Getenv("PIIO_ID")

	if httpPort == "" {
		httpPort = "8080"
	}

	db, err := psql.New(dbURL)
	if err != nil {
		log.Fatalf("unable to connect to db %s", err)
	}

	pool := &worker.WorkerPool{
		Database: db,
	}

	imgResizer := &resizer.Worker{
		WorkerPool:  pool,
		AWSBucket:   awsBucket,
		BlitlineID:  blitlineID,
		CallbackURL: blitlineCallbackURL,
		Poll:        env == Development,
	}

	err = imgResizer.Register()
	if err != nil {
		log.Fatalf("unable to register AWS S3 job %s", err)
	}

	go pool.Start()

	ub, err := urlbuilder.New(hashidSalt)
	if err != nil {
		log.Fatalf("unable to initialize URL builder %s", err)
	}

	gen := generator.Generator{
		Database: db,
	}

	auth := authentication.Authenticator{}

	session := &session.Manager{
		Database: db,
		Secret:   sessionSecret,
	}

	tpl := html.New("web/templates", piioDomain, piioID)

	srv := &server.Server{
		Template:          tpl,
		Database:          db,
		URLBuilder:        ub,
		PracticeGenerator: gen,
		Authenticator:     auth,
		SessionManager:    session,
		ImageResizer:      imgResizer,
	}

	mux := srv.NewServeMux()

	fmt.Printf("server listening on port %s\n", httpPort)

	err = http.ListenAndServe(":"+httpPort, mux)
	if err != nil {
		log.Fatalf("unable to start server %s", err)
	}
}
