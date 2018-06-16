.PHONY: db

docker:
	docker build --rm -t gcr.io/cyberbrain-app/server .

push:
	gcloud container builds submit --tag gcr.io/cyberbrain-app/server .
