FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /app/examples/kafka_standalone/topic/delete
RUN go build -o main -tags musl delete.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/examples/kafka_standalone/topic/delete/main .
CMD ./main