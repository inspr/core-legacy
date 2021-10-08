FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /app/examples/kafka_standalone/reader
RUN go build -o main -tags musl reader.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/examples/kafka_standalone/reader/main .
CMD ./main