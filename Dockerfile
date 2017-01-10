FROM alpine:3.5

RUN apk add --no-cache --update ca-certificates

COPY bin/k8sec /k8sec

ENTRYPOINT ["/k8sec"]
