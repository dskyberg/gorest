FROM docker.dev.confyrm.com/golang:1

ENV GOREST_REPO github.com/confyrm/gorest

RUN go get $GOREST_REPO && go install $GOREST_REPO \
  && ln -sf $GOPATH/bin/gorest /data/bin \
  && ln -sf $GOPATH/src/$GOREST_REPO/help /data/conf \
  && apk --purge del build-dependencies && rm -rf /var/cache/apk/* && rm -rf $GOPATH/src/ $GOPATH/bin/godep $GOPATH/pkg/

ENTRYPOINT ["/gopath/bin/gorest"]
