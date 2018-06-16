dev:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/server cmd/server/main.go
	docker build --rm -f Dockerfile.dev -t cyberbrain/server:latest .

docker:
	docker build --rm -t cyberbrain/server:latest .
	docker push cyberbrain/server:latest
