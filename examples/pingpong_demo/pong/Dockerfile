FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build examples/pingpong_demo/pong/pong.go

FROM alpine
WORKDIR /app
CMD ./pong