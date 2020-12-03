FROM golang:1.15.5 AS build-env
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

FROM build-env AS build
RUN go build -o k8s -tags musl cmd/operator/k8s/main.go 

FROM alpine AS deploy
COPY --from=build /app/k8s /k8s
ENTRYPOINT ["/k8s"]
