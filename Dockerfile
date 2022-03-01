FROM golang:1.17 AS builder

WORKDIR /go/src/github.com/dtan4/k8sec

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /k8sec

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /k8sec /k8sec

ENTRYPOINT ["/k8sec"]
