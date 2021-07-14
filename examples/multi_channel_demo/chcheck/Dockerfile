FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o main examples/multi_channel_demo/chcheck/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/main .
CMD ./main
