FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/examples/primes/filter
RUN go build

FROM alpine
WORKDIR /app
COPY --from=build /app/examples/primes/filter/filter filter
CMD ./filter

