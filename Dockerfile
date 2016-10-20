FROM alpine:3.4

ENV K8SEC_VERSION v0.1.1

RUN apk add --no-cache --update ca-certificates unzip wget \
    && wget -qO /k8sec.zip "https://github.com/dtan4/k8sec/releases/download/${K8SEC_VERSION}/k8sec-${K8SEC_VERSION}-linux-amd64.zip" \
    && unzip /k8sec.zip -d /bin \
    && apk del --purge unzip wget \
    && rm -rf /k8sec.zip

ENTRYPOINT ["/bin/linux-amd64/k8sec"]
