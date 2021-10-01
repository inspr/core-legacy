FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o main examples/route_demo/api/api.go

FROM alpine AS final
WORKDIR /app
COPY --from=build /app/main .
CMD ./main