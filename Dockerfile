FROM alpine:3.4

ENV K8SEC_VERSION v0.1.0
ENV GLIBC_VERSION 2.23-r3

RUN apk add --no-cache --update ca-certificates unzip wget \
    && wget -qO /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
    && wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk \
    && apk add --no-cache glibc-${GLIBC_VERSION}.apk \
    && wget -qO /k8sec.zip "https://github.com/dtan4/k8sec/releases/download/${K8SEC_VERSION}/k8sec-${K8SEC_VERSION}-linux-amd64.zip" \
    && unzip /k8sec.zip -d /bin \
    && apk del --purge unzip wget \
    && rm -rf glibc-${GLIBC_VERSION}.apk /k8sec.zip

ENTRYPOINT ["/bin/linux-amd64/k8sec"]
