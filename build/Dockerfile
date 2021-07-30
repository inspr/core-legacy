FROM golang:alpine AS build

RUN apk update && apk add make git gcc musl-dev

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY Makefile .
ENV GOFLAGS -tags=musl
RUN make go/download
COPY . .
RUN make go/build

FROM alpine
WORKDIR /deploy

RUN addgroup inspr
RUN adduser inspr -G inspr -D

USER inspr
ARG TARGET

COPY --from=build /build/build/bin/${TARGET} ./entrypoint
ENTRYPOINT ./entrypoint
