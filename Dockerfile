### source build ###
FROM golang:1.13-alpine3.10 as build

COPY . /src

WORKDIR /src

RUN set -ex ;\
    apk add git ;\
    go get -d -v -t ;\
    CGO_ENABLED=0 GOOS=linux go build -v -o /files/usr/local/bin/request-forwarder

### runtime build ###
FROM centos:7

COPY --from=build /files /

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/request-forwarder" ]