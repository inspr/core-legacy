# builds kafka sidecar
FROM gcr.io/red-inspr/inspr
WORKDIR /inspr/cmd/sidecars
RUN go build kafka -tags musl