FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o main examples/pingpong_demo/ping/ping.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/main .
CMD ./main