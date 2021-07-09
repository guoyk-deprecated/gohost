FROM golang:1.16 AS build
ENV PROJECT gohost
ENV GOPROXY https://goproxy.io
WORKDIR /src/$PROJECT
COPY . .
RUN CGO_ENABLED=0 go install -mod vendor -a -tags netgo -ldflags=-w

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
WORKDIR /opt/gohost
COPY --from=build /go/bin/gohost /opt/gohost/gohost
ADD views /opt/gohost/views
CMD [ "/opt/gohost/gohost" ]
