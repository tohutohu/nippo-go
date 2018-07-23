FROM golang:1.10.2-alpine3.7
RUN apk update \
  && apk add --no-cache git \
  && go get -u github.com/golang/dep/cmd/dep 
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

WORKDIR /go/src/github.com/tohutohu/nippo
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only=true

ENV GOPATH=/go
COPY . .
RUN go build -o app
CMD ./app
