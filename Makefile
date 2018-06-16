dev:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/server cmd/server/main.go
	docker build --rm -f Dockerfile.dev -t cyberbrain/server:latest .

release:
	docker build --rm --no-cache  -f Dockerfile -t cyberbrain/server:v${version} .
	docker push cyberbrain/server:v${version}
