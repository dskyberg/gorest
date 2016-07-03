FROM docker.dev.confyrm.com/golang:1

ENV GOREST_REPO github.com/confyrm/gorest

RUN go get -d -u $GOREST_REPO && go install $GOREST_REPO && \
    apk --purge del build-dependencies && rm -rf /var/cache/apk/* && rm -rf $GOPATH/src/ $GOPATH/bin/godep $GOPATH/pkg/

ENTRYPOINT ["/gopath/bin/gorest"]
