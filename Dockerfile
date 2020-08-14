FROM golang:1.14 as build
LABEL maintainer="Ruihao Liu <liuruihao@huobi.com>"
COPY . /tmp/tcpproxy

RUN cd /tmp/tcpproxy && go get . && go build

FROM ubuntu:18.04
LABEL maintainer="Ruihao Liu <liuruihao@huobi.com>"
COPY --from=build /tmp/tcpproxy/stratumtcpproxy /usr/local/bin/proxy
ENTRYPOINT ["proxy"]