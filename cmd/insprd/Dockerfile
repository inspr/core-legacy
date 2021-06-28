FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build-env

WORKDIR /inspr
COPY go.mod go.mod
RUN go mod download
COPY . .

WORKDIR /inspr/cmd/insprd
RUN go build -tags musl

FROM gcr.io/insprlabs/inspr/sidecar/kafka:deploy
WORKDIR /deploy
COPY --from=build-env /inspr/cmd/insprd/insprd insprd
ENTRYPOINT ./insprd
