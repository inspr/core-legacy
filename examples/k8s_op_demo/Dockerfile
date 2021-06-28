FROM gcr.io/insprlabs/inspr/sidecar/kafka:build AS build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/cmd/insprd/operators/kafka/nodes/test/
RUN go build -tags musl


FROM gcr.io/insprlabs/inspr/sidecar/kafka:deploy

WORKDIR /app
COPY --from=build /app/cmd/insprd/operators/kafka/nodes/test/test test
CMD ./test

