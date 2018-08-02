FROM golang:alpine

RUN apk update && \
    apk add git

WORKDIR /go/src/github.com/ttpham0111/exp-ose/exp

COPY exp .

RUN go get -v && \
    go install -v

ENTRYPOINT ["exp"]
