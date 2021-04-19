FROM golang:1.16.3-alpine3.13 AS build

WORKDIR /inspr
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/uid_provider
RUN go build

FROM alpine
WORKDIR /deploy
COPY --from=build /inspr/cmd/uid_provider/uid_provider uid_provider
ENTRYPOINT ./uid_provider
