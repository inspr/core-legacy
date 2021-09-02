FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o main examples/bench_mark/receiver/receiver.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/main .
CMD ./main