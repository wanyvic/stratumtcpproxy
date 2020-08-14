FROM golang:1.14 as build
LABEL maintainer="Ruihao Liu <liuruihao@huobi.com>"

COPY . /tmp/tcpproxy

RUN cd /tmp/tcpproxy && go get . && go build

FROM ubuntu:18.04
LABEL maintainer="Ruihao Liu <liuruihao@huobi.com>"
ARG APT_MIRROR_URL

COPY update_apt_sources.sh /usr/local/bin/
RUN /usr/local/bin/update_apt_sources.sh "$APT_MIRROR_URL"

# add DNS
RUN echo nameserver 114.114.114.114 >> /etc/resolv.conf && \
    echo nameserver 223.5.5.5       >> /etc/resolv.conf && \
    echo nameserver 8.8.8.8         >> /etc/resolv.conf


COPY --from=build /tmp/tcpproxy/stratumtcpproxy /usr/local/bin/proxy

COPY ./entrypoint.sh ./wait-for-it.sh /

# entrypoint
ENTRYPOINT ["/entrypoint.sh"]
