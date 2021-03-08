FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/examples/printer
RUN go build -tags musl


FROM alpine
WORKDIR /app
COPY --from=build /app/examples/printer printer
CMD ./printer