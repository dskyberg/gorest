FROM docker.dev.confyrm.com/golang:2

ENV GOREST_REPO github.com/confyrm/gorest

RUN go get $GOREST_REPO && go install $GOREST_REPO \
    && apk --purge del build-dependencies && rm -rf /var/cache/apk/*

RUN mkdir /data/logs

CMD ["/gopath/bin/gorest"]
