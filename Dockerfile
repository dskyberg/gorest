FROM docker.dev.confyrm.com/golang:2

ENV GOREST_REPO github.com/confyrm/gorest
ENV APP_ROOT $GOPATH/src/$GOREST_REPO

RUN go get $GOREST_REPO && go install $GOREST_REPO \
    && apk --purge del build-dependencies && rm -rf /var/cache/apk/*

CMD ["/gopath/bin/gorest"]
