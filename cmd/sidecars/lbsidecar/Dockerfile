FROM golang:alpine as build
WORKDIR /inspr
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . . 
WORKDIR /inspr/cmd/sidecars/lbsidecar
RUN go build


FROM alpine
WORKDIR /deploy
COPY --from=build /inspr/cmd/sidecars/lbsidecar/lbsidecar lbsidecar
ENTRYPOINT ./lbsidecar
