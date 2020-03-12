FROM gcr.io/distroless/static

COPY k8sec /

ENTRYPOINT ["/k8sec"]
