FROM docker.dev.confyrm.com/golang:2

ENV GOREST_REPO github.com/confyrm/gorest
ENV GOREST_ROOT $GOPATH/src/$GOREST_REPO
ENV APP_ROOT /gorest/config
ENV GOREST_CONFIG $APP_ROOT/config.yml

RUN mkdir -p $APP_ROOT

# Install project dependencies
RUN go get -v $GOREST_REPO && rm -Rf $GOREST_ROOT \
  && rm -f $GOPATH/bin/gorest

# Ensure docker is not cached the wrong codebase state
ADD . $GOREST_ROOT
COPY  dist/ $APP_ROOT/
RUN ls $APP_ROOT

RUN cd $GOREST_ROOT \
  && go get \
  && go install \
  && cd \
  && rm -Rf $GOREST_ROOT

CMD ["/gopath/bin/gorest"]
