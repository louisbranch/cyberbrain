FROM golang:alpine AS build-env
ENV repo gitlab.com/luizbranco/cyberbrain
WORKDIR /app
RUN apk update && apk add git && apk add ca-certificates
RUN go get -u ${repo}/cmd/server
RUN go build -o /app/server ${repo}/cmd/server

FROM alpine
ENV src /go/src/gitlab.com/luizbranco/cyberbrain
WORKDIR /app
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env ${src}/web/assets /app/web/assets
COPY --from=build-env ${src}/web/templates /app/web/templates
COPY --from=build-env /app/server /app

EXPOSE 8080
ENTRYPOINT ./server
