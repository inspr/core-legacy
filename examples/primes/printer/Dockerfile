FROM golang:alpine as build
WORKDIR /app
COPY . . 
RUN go build -o main examples/primes/printer/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/main .
CMD ./main