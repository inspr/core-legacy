FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/examples/filter
RUN go build
# RUN chmod +x /app/examples/generator/generator

FROM alpine
WORKDIR /app
COPY --from=build /app/examples/filter filter
CMD ./filter

