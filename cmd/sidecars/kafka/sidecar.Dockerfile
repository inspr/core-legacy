FROM gcr.io/red-inspr/inspr AS build-env

# install go lib dependencies
WORKDIR /inspr
COPY go.mod go.mod
RUN go mod download
COPY . .

# builds kafka sidecar
WORKDIR /inspr/cmd/sidecars/kafka
RUN go build -tags musl