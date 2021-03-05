FROM gcr.io/red-inspr/inspr/examples/generator:build AS build
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . . 
WORKDIR /app/examples/num_generator
RUN go build -tags musl


FROM gcr.io/red-inspr/inspr/examples/generator:deploy
WORKDIR /app
COPY --from=build /app/examples/num_generator num_generator
CMD ./num_generator

