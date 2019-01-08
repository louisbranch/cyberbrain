package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/luizbranco/cyberbrain/authentication"
	"gitlab.com/luizbranco/cyberbrain/db/psql"
	"gitlab.com/luizbranco/cyberbrain/s3img"
	"gitlab.com/luizbranco/cyberbrain/web/html"
	"gitlab.com/luizbranco/cyberbrain/web/server"
	"gitlab.com/luizbranco/cyberbrain/web/session"
	"gitlab.com/luizbranco/cyberbrain/web/urlbuilder"
	"gitlab.com/luizbranco/cyberbrain/worker"
	"gitlab.com/luizbranco/cyberbrain/worker/offline"
	"gitlab.com/luizbranco/cyberbrain/worker/resizer"
)

const (
	Development = "dev"
	Production  = "prod"
)

func main() {
	httpPort := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	sessionSecret := os.Getenv("SESSION_SECRET")
	hashidSalt := os.Getenv("HASHID_SALT")
	env := os.Getenv("ENVIRONMENT")

	if env == "" {
		env = Development
	}

	blitlineID := os.Getenv("BLITLINE_ID")
	blitlineCallbackURL := os.Getenv("BLITLINE_CALLBACK_URL")

	awsID := os.Getenv("AWS_ID")
	awsSecret := os.Getenv("AWS_SECRET")
	awsRegion := os.Getenv("AWS_REGION")
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

	var imgResizer worker.ImageResizer

	if env == Production {
		blitlineResizer := &resizer.Worker{
			WorkerPool:  pool,
			AWSBucket:   awsBucket,
			BlitlineID:  blitlineID,
			CallbackURL: blitlineCallbackURL,
			Poll:        env == Development,
		}

		err := blitlineResizer.Register()
		if err != nil {
			log.Fatalf("unable to register AWS S3 job %s", err)
		}

		imgResizer = blitlineResizer
	} else {
		imgResizer = &offline.ImageOfflineResizer{}
	}

	go pool.Start()

	ub, err := urlbuilder.New(hashidSalt)
	if err != nil {
		log.Fatalf("unable to initialize URL builder %s", err)
	}

	auth := authentication.Authenticator{}

	session := &session.Manager{
		Database: db,
		Secret:   sessionSecret,
	}

	tpl := html.New("web/templates", env, piioDomain, piioID)

	imgUploader, err := s3img.New(awsID, awsSecret, awsRegion, awsBucket)
	if err != nil {
		log.Fatalf("unable to initialize S3 image uploader %s", err)
	}

	srv := &server.Server{
		Template:       tpl,
		Database:       db,
		URLBuilder:     ub,
		Authenticator:  auth,
		SessionManager: session,
		ImageResizer:   imgResizer,
		ImageUploader:  imgUploader,
	}

	mux := srv.NewServeMux()

	fmt.Printf("server listening on port %s\n", httpPort)

	err = http.ListenAndServe(":"+httpPort, mux)
	if err != nil {
		log.Fatalf("unable to start server %s", err)
	}
}
