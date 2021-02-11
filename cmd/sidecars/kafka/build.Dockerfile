# librdkafka base image
FROM golang:alpine AS kafka-build 

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git gcc g++ make bash pkgconfig

RUN git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install