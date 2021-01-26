FROM golang:alpine AS build-env

WORKDIR /inspr
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/insprd
RUN go build

FROM alpine
WORKDIR /deploy
COPY --from=build-env /inspr/cmd/insprd/insprd insprd
ENTRYPOINT ./insprd