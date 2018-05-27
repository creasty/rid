FROM golang:1.10-alpine

#  Dep
#-----------------------------------------------
ENV DEP_VERSION 0.4.1

RUN set -eux \
  curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -o $GOPATH/bin/dep \
  && chmod +x $GOPATH/bin/dep

#  Golint
#-----------------------------------------------
RUN go get -u github.com/golang/lint/golint

#  Library
#-----------------------------------------------
RUN apk add --no-cache -U \
    bash \
    build-base \
    coreutils \
    make

#  App
#-----------------------------------------------
ENV APP_DIR
WORKDIR $APP_DIR
RUN ln -sf $APP_DIR /app
