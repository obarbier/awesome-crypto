FROM golang:1.17-alpine
ENV GO111MODULE=on

RUN apk update && apk upgrade && \
    apk add --no-cache bash git curl

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY config.yml config.yml
COPY internal internal
COPY pkg/oauth pkg/oauth
COPY cmd/oauth cmd/oauth

# building servic
 RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/oauth/

# command to run at the immediate start of the container
ENTRYPOINT ["./oauth", "serve"]
