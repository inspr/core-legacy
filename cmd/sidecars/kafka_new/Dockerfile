# builds kafka sidecar
FROM gcr.io/red-inspr/inspr/sidecar/kafka:build AS build-env

WORKDIR /inspr
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/sidecars/kafka_new
RUN go build -o kafka_new -tags musl

FROM gcr.io/red-inspr/inspr/sidecar/kafka:deploy
WORKDIR /deploy
COPY --from=build-env /inspr/cmd/sidecars/kafka_new/kafka_new kafka_new
ENTRYPOINT ./kafka_new