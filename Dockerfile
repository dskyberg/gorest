FROM docker.dev.confyrm.com/alpine:3

ENV GOREST_REPO github.com/confyrm/gorest

ENV GOROOT /usr/lib/go
ENV GOPATH /gopath
ENV GOBIN /gopath/bin
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN apk add -U -t build-dependencies go git make && \
    go get -d -u $GOREST_REPO && cd $GOPATH/src/$GOREST_REPO && go install && \
    apk --purge del build-dependencies && rm -rf /var/cache/apk/* && rm -rf $GOPATH/src/ $GOPATH/bin/godep $GOPATH/pkg/


ENTRYPOINT ["/gopath/bin/gorest"]
