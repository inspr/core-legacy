# Install librdkafka base image, which is necessary for kafka to run
# Tag used to build this image:
# gcr.io/insprlabs/inspr/sidecar/kafka:deploy
FROM alpine AS kafka-deploy

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git gcc g++ make bash pkgconfig

RUN git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install