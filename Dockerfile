FROM alpine:3.10

RUN apk add --no-cache ca-certificates

WORKDIR /app

ADD views /app/views

ADD gohost /

CMD ["/gohost"]
