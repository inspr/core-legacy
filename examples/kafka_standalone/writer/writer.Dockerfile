FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /app/examples/kafka_standalone/writer
RUN go build -o main -tags musl writer.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/examples/kafka_standalone/writer/main .
CMD ./main