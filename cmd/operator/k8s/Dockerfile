FROM gcr.io/red-inspr/inspr AS build-env
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .

FROM build-env AS build
RUN go build -o k8s -tags musl main.go 

FROM alpine AS deploy
COPY --from=build /usr/lib/*kafka* /usr/lib/
COPY --from=build /app/k8s /k8s
ENTRYPOINT ["/k8s"]
