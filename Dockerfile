FROM golang:1.14 AS builder

WORKDIR /go/src/github.com/dtan4/k8sec
COPY . /go/src/github.com/dtan4/k8sec

RUN CGO_ENABLED=0 make

FROM alpine:3.8

RUN apk add --no-cache --update ca-certificates

COPY --from=builder /go/src/github.com/dtan4/k8sec/bin/k8sec /k8sec

ENTRYPOINT ["/k8sec"]
