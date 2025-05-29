# k8sec

[![GitHub Actions](https://github.com/dtan4/k8sec/workflows/Test/badge.svg)](https://github.com/dtan4/k8sec/actions?query=workflow%3ATest+branch%3Amaster)
[![GitHub release](https://img.shields.io/github/release/dtan4/k8sec.svg)](https://github.com/dtan4/k8sec/releases)

CLI tool to manage [Kubernetes Secrets](http://kubernetes.io/docs/user-guide/secrets/) easily

## Requirements

Kubernetes 1.18 or above

## Installation

### Using Homebrew (OS X only)

Formula is available at [dtan4/homebrew-dtan4](https://github.com/dtan4/homebrew-tools).

```bash
brew tap dtan4/tools
brew install k8sec
```

### Precompiled binary

Precompiled binaries for Windows, OS X, Linux are available at [Releases](https://github.com/dtan4/k8sec/releases).

### From source

```sh-session
go get -d github.com/dtan4/k8sec
cd $GOPATH/src/github.com/dtan4/k8sec
make install
```

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

```sh-session
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

```sh-session
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

```sh-session
$ k8sec unset NAME KEY1 KEY2...

# Example
$ k8sec unset rails rails-env
```

### `k8sec load`

Load secrets from dotenv (key=value) format text

```sh-session
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

```sh-session
$ k8sec dump [-f FILENAME] [--noquotes] [NAME]

# Example
$ k8sec dump rails
database-url="postgres://example.com:5432/dbname"

# Save as .env
$ k8sec dump -f .env rails
$ cat .env
database-url="postgres://example.com:5432/dbname"

# Save as .env without qoutes
$ k8sec dump -f .env --noquotes rails
$ cat .env
database-url=postgres://example.com:5432/dbname
```

## Contribution

1. Fork ([https://github.com/dtan4/k8sec/fork](https://github.com/dtan4/k8sec/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[dtan4](https://github.com/dtan4)

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
