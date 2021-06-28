FROM golang:alpine AS build

WORKDIR /inspr
COPY go.mod go.mod
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/authsvc
RUN go build -tags musl

FROM alpine
WORKDIR /deploy
COPY --from=build /inspr/cmd/authsvc/authsvc authsvc
ENTRYPOINT ./authsvc
