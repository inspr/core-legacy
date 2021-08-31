FROM golang:alpine AS build

WORKDIR /inspr
COPY go.mod go.mod
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/authsvc/secret
RUN go build

FROM alpine
WORKDIR /deploy
COPY --from=build /inspr/cmd/authsvc/secret/secret secret
ENTRYPOINT ./secret
