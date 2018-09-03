start:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/server cmd/server/main.go
	docker build --rm -f Dockerfile.dev -t cyberbrain/server:dev .
	docker-compose up -d

stop:
	docker-compose stop

ci:
	go vet ./...
	go test ./...

release:
	git tag ${version}
	git push origin ${version}
	docker build --rm --no-cache  -f Dockerfile -t cyberbrain/server:${version} .
	docker push cyberbrain/server:${version}

local:
	DATABASE_URL="postgres://cyberbrain:s3cr3t@localhost:5432/cyberbrain?sslmode=disable" SESSION_SECRET=S5L56UH4TKQNYJMNK486WLH7E4MQRV26 HASHID_SALT=LXRQSJ68ZUXMEQOP2E0T7QJW43GGC3FA go run cmd/server/main.go
