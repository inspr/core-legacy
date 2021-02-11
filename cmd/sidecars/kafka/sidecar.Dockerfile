# builds kafka sidecar
FROM gcr.io/red-inspr/inspr
WORKDIR /inspr/cmd/sidecars/kafka
RUN go build -tags musl