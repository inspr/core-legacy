FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build examples/pingpong_demo/ping/ping.go

FROM alpine
WORKDIR /app
CMD ./ping