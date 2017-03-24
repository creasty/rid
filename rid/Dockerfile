FROM golang:1.7.5-alpine


#  Glide
#-----------------------------------------------
ENV GLIDE_VERSION 0.11.1

RUN apk add --no-cache -U --virtual .build-deps \
    curl \
  && curl -fL "https://github.com/Masterminds/glide/releases/download/v$GLIDE_VERSION/glide-v$GLIDE_VERSION-linux-amd64.zip" -o glide.zip \
  && unzip glide.zip \
  && mv ./linux-amd64/glide /usr/local/bin/glide \
  && rm -fr ./linux-amd64 \
  && rm ./glide.zip \
  && apk del .build-deps \
  \
  && apk add --no-cache -U --virtual .glide-deps \
    git


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
ENV APP_DIR /go/src/github.com/creasty/rid
WORKDIR $APP_DIR
RUN ln -sf $APP_DIR /app
