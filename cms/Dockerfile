FROM alpine:3.13

WORKDIR /root/

COPY cms cms
COPY assets assets
COPY env/sample.config env/config

ENTRYPOINT ["./cms"]
