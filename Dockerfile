FROM golang:alpine AS build-env
ENV src /go/src/gitlab.com/luizbranco/cyberbrain
WORKDIR $src
ADD . $src
RUN apk add --no-cache git
RUN cd ${src}/cmd/server && go get -u && go build -o server

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENV src /go/src/gitlab.com/luizbranco/cyberbrain
WORKDIR /app
COPY --from=build-env ${src}/web/assets /app/web/assets
COPY --from=build-env ${src}/web/templates /app/web/templates
COPY --from=build-env ${src}/cmd/server/server /app

EXPOSE 8080
ENTRYPOINT ./server
