package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/endpoints"
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

var httpPort, dbURL, sessionSecret, hashidSalt, awsID, awsSecret, awsBucket,
	awsRegion string

func init() {
	httpPort = os.Getenv("HTTP_PORT")
	dbURL = os.Getenv("DATABASE_URL")
	sessionSecret = os.Getenv("SESSION_SECRET")
	hashidSalt = os.Getenv("HASHID_SALT")

	if httpPort == "" {
		httpPort = "8080"
	}

	if dbURL == "" {
		dbURL = "postgres://cyberbrain:s3cr3t@192.168.0.11:5432/cyberbrain?sslmode=disable"
	}

	if sessionSecret == "" {
		sessionSecret = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	}

	if hashidSalt == "" {
		hashidSalt = "s3cret"
	}

	flag.StringVar(&httpPort, "http-port", httpPort, "http port to start server")
	flag.StringVar(&dbURL, "database-url", dbURL, "database connection url")
	flag.StringVar(&sessionSecret, "session-secret", sessionSecret, "session cookie id secret")
	flag.StringVar(&hashidSalt, "hashid-salt", hashidSalt, "salt for hashid url")
	flag.StringVar(&awsID, "aws-id", "", "AWS Access Key ID")
	flag.StringVar(&awsSecret, "aws-secret", "", "AWS Access Key Secret")
	flag.StringVar(&awsBucket, "aws-bucket", "", "AWS Bucket Name")
	flag.StringVar(&awsRegion, "aws-region", endpoints.UsEast2RegionID, "AWS Region")

	flag.Parse()
}

func main() {
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
		log.Fatalf("unable to initialize urbuilder %s", err)
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
