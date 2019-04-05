# k8sec

[![Build Status](https://travis-ci.org/wantedly/k8sec.svg?branch=master)](https://travis-ci.org/wantedly/k8sec)
[![GitHub release](https://img.shields.io/github/release/wantedly/k8sec.svg)](https://github.com/wantedly/k8sec/releases)
[![Docker Repository on Quay](https://quay.io/repository/wantedly/k8sec/status "Docker Repository on Quay")](https://quay.io/repository/wantedly/k8sec)

CLI tool to manage [Kubernetes Secrets](http://kubernetes.io/docs/user-guide/secrets/) easily

## Requirements

Kubernetes 1.3 or above

## Installation

### Precompiled binary

Precompiled binaries for Windows, OS X, Linux are available at [Releases](https://github.com/wantedly/k8sec/releases).

### From source

```bash
$ go get -d github.com/wantedly/k8sec
$ cd $GOPATH/src/github.com/wantedly/k8sec
$ make deps
$ make install
```

### Docker image

Docker image is available at [`quay.io/wantedly/k8sec`](https://quay.io/repository/wantedly/k8sec).

## Usage

### Global options

|Option|Description|Required|Default|
|---------|-----------|-------|-------|
|`--context=CONTEXT`|Kubernetes context|||
|`--kubeconfig=KUBECONFIG`|Path of kubeconfig||`~/.kube/config`|
|`-n`, `--namespace=NAMESPACE`|Kubernetes namespace||`default`|
|`-h`, `-help`|Print command line usage|||

### `k8sec list`

List secrets

```bash
$ k8sec list [--base64] [NAME]

# Example
$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"

# Show values as base64-encoded string
$ k8sec list --base64 rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    cG9zdGdyZXM6Ly9leGFtcGxlLmNvbTo1NDMyL2RibmFtZQ==
```

### `k8sec set`

Set secrets

```bash
$ k8sec set [--base64] NAME KEY1=VALUE1 [KEY2=VALUE2 ...]

$ k8sec set rails rails-env=production
rails

# Set base64-encoded value
$ echo -n dtan4 | base64
ZHRhbjQ=
$ k8sec set --base64 rails foo=ZHRhbjQ=
rails

# Result
$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"
rails   Opaque  foo             "dtan4"
```

### `k8sec unset`

Unset secrets

```bash
$ k8sec unset NAME KEY1 KEY2...

# Example
$ k8sec unset rails rails-env
```

### `k8sec load`

Load secrets from dotenv (key=value) format text

```bash
$ k8sec load [-f FILENAME] NAME

# Example
$ cat .env
database-url="postgres://example.com:5432/dbname"
$ k8sec load -f .env rails

# Load from stdin
$ cat .env | k8sec load rails
```

### `k8sec dump`

Dump secrets as dotenv (key=value) format

```bash
$ k8sec dump [-f FILENAME] [NAME]

# Example
$ k8sec dump rails
database-url="postgres://example.com:5432/dbname"

# Save as .env
$ k8sec dump -f .env rails
$ cat .env
database-url="postgres://example.com:5432/dbname"
```

## Contribution

Go 1.8 or above is required.

1. Fork ([https://github.com/wantedly/k8sec/fork](https://github.com/wantedly/k8sec/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[dtan4](https://github.com/dtan4)
[wantedly](https://github.com/wantedly)

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
