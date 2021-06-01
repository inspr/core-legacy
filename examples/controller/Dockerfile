FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .
WORKDIR /app/examples/controller
RUN go build

FROM alpine
WORKDIR /app
COPY --from=build /app/examples/controller/controller controller
CMD ./controller

