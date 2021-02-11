FROM gcr.io/red-inspr/inspr AS build-env

WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

FROM alpine AS deploy
WORKDIR /app
COPY --from=build-env /inspr/cmd/sidecars/kafka /app/kafka