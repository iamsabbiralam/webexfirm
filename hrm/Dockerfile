FROM alpine:3.13

WORKDIR /root/

COPY hrm hrm
ADD env/sample.config env/config
ADD migrations/sql migrations/sql

ENTRYPOINT ["./hrm"]