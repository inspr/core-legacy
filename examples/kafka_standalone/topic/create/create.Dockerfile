FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /app/examples/kafka_standalone/topic/create
RUN go build -o main -tags musl create.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/examples/kafka_standalone/topic/create/main .
CMD ./main