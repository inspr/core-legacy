FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/examples/generator
RUN go build


FROM alpine
WORKDIR /app
COPY --from=build /app/examples/generator/generator generator
CMD ./generator

