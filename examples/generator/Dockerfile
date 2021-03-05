FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/examples/generator
RUN go build -tags musl


FROM scratch
WORKDIR /app
COPY --from=build /app/examples/generator generator
CMD ./generator

