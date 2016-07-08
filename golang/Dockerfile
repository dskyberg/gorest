FROM docker.dev.confyrm.com/alpine:3

ENV GOROOT /usr/lib/go
ENV GOPATH /gopath
ENV GOBIN $GOPATH/bin
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN echo http://dl-4.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories \
  && apk update && apk add libgit2 \
  && apk add -U -t build-dependencies go git \
    && rm -rf /var/cache/apk/*
