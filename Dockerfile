FROM golang:1.10.2-alpine3.7
RUN apk update \
  && apk add --no-cache git \
  && go get -u github.com/golang/dep/cmd/dep 

ENV GOPATH=/go
WORKDIR /go/src/github.com/tohutohu/nippo
COPY . /go/src/github.com/tohutohu/nippo
RUN dep ensure && go build -o app
CMD /go/src/github.com/tohutohu/nippo/app
