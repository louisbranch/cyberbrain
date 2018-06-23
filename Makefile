dev:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/server cmd/server/main.go
	docker build --rm -f Dockerfile.dev -t cyberbrain/server:dev .
	manifold run -- docker-compose up

ci:
	go test ./...

release:
	git tag ${version}
	git push origin ${version}
	docker build --rm --no-cache  -f Dockerfile -t cyberbrain/server:${version} .
	docker push cyberbrain/server:${version}
