FROM docker.dev.confyrm.com/golang:2

ENV GOREST_REPO github.com/confyrm/gorest
ENV APP_ROOT $GOPATH/src/$GOREST_REPO

# Install project dependencies
RUN go get -v $GOREST_REPO && rm -Rf $APP_ROOT

# Ensure docker is not cached the wrong codebase state
ADD . $APP_ROOT

RUN go install $GOREST_REPO

CMD ["/gopath/bin/gorest"]
