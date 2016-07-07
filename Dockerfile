FROM docker.dev.confyrm.com/golang:1

ENV GOREST_REPO github.com/confyrm/gorest

RUN go get $GOREST_REPO && go install $GOREST_REPO \
  && ln -sf /data/bin $GOPATH/bin/gorest \
  && ln -sf /data/conf $GOPATH/src/$GOREST_REPO/help \
  && apk --purge del build-dependencies && rm -rf /var/cache/apk/* && rm -rf $GOPATH/src/ $GOPATH/bin/godep $GOPATH/pkg/

ENTRYPOINT ["/gopath/bin/gorest"]
