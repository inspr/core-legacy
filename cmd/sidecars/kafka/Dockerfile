# builds kafka sidecar
FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build-env

WORKDIR /inspr
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/sidecars/kafka
RUN go build -o kafka -tags musl

FROM gcr.io/insprlabs/inspr/sidecar/kafka:deploy
WORKDIR /deploy
COPY --from=build-env /inspr/cmd/sidecars/kafka/kafka kafka
ENTRYPOINT ./kafka