FROM alpine:3.21 AS builder

ARG LIBCLERI_VERSION=1.0.2
ARG SIRIDB_VERSION=2.0.52

RUN apk update \
 && apk upgrade \
 && apk add --no-cache --update \
        gcc \
        git \
        libuv-dev \
        linux-headers \
        make \
        musl-dev \
        pcre2-dev \
        util-linux-dev \
        yajl-dev \
 && git clone https://github.com/cesbit/libcleri.git \
              /tmp/libcleri \
 && cd /tmp/libcleri \
 && cd Release \
 && make all \
 && make install \
 && git clone https://github.com/SiriDB/siridb-server.git \
        /tmp/siridb-server \
 && cd /tmp/siridb-server \
 && git checkout $SIRIDB_VERSION \
 && cd Release \
 && make clean \
 && make

 FROM alpine:3.21 AS stage

 COPY --from=builder /tmp/siridb-server/Release/siridb-server /usr/local/bin/
 COPY --from=builder /usr/lib/libcleri* /usr/lib/

 RUN apk add --no-cache --update \
        dumb-init \
        libuuid \
        libuv \
        pcre2 \
        su-exec \
        yajl \
 && mkdir -p /etc/siridb \
 && mkdir -p /var/lib/siridb \
 && chmod +x /usr/local/bin/*

EXPOSE 9000
EXPOSE 9010
EXPOSE 9200
EXPOSE 9020
 CMD ["/usr/local/bin/siridb-server", "--log-level", "debug"]