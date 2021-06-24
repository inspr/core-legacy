FROM golang:alpine AS build-env
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
WORKDIR /app/examples/multi_channel_demo/chtest
RUN go build

FROM alpine
WORKDIR /app
COPY --from=build-env /app/examples/multi_channel_demo/chtest/chtest .
CMD ./chtest
