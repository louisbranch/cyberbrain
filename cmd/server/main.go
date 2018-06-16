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
	"gitlab.com/luizbranco/cyberbrain/worker/jobs/s3img"
)

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	dbURL := os.Getenv("DATABASE_URL")
	sessionSecret := os.Getenv("SESSION_SECRET")
	hashidSalt := os.Getenv("HASHID_SALT")

	awsID := os.Getenv("AWS_ID")
	awsSecret := os.Getenv("AWS_SECRET")
	awsBucket := os.Getenv("AWS_BUCKET")
	awsRegion := os.Getenv("AWS_REGION")

	if httpPort == "" {
		httpPort = "8080"
	}

	if sessionSecret == "" {
		sessionSecret = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	}

	if hashidSalt == "" {
		hashidSalt = "s3cret"
	}

	db, err := psql.New(dbURL)
	if err != nil {
		log.Fatalf("unable to connect to db %s", err)
	}

	pool := &worker.WorkerPool{
		Database: db,
	}

	s3 := &s3img.Worker{
		Database:   db,
		WorkerPool: pool,
		AWSID:      awsID,
		AWSSecret:  awsSecret,
		AWSBucket:  awsBucket,
		AWSRegion:  awsRegion,
	}

	err = s3.Register()
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

	tpl := html.New("web/templates")

	srv := &server.Server{
		Template:          tpl,
		Database:          db,
		URLBuilder:        ub,
		PracticeGenerator: gen,
		Authenticator:     auth,
		SessionManager:    session,
		ImageUploader:     s3,
	}

	mux := srv.NewServeMux()

	fmt.Printf("Server listening on port %s\n", httpPort)

	err = http.ListenAndServe(":"+httpPort, mux)
	if err != nil {
		log.Fatalf("unable to start server %s", err)
	}
}
