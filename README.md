# k8sec

CLI tool to manage [Kubernetes Secrets](http://kubernetes.io/docs/user-guide/secrets/) easily.

## Usage

### `k8sec list`

List secrets

``` bash
$ k8sec list [--base64] [--kubeconfig KUBECONFIG] [--namespace NAMESPACE] [NAME]

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

``` bash
$ k8sec set [--base64] [--kubeconfig KUBECONFIG] [--namespace NAMESPACE] NAME KEY1=VALUE1 KEY2=VALUE2

# Example
$ k8sec set rails rails-env=production
rails

# Pass base64-encoded value
$ echo dtan4 | base64
ZHRhbjQK
$ k8sec set --base64 rails foo=ZHRhbjQK
rails
$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"
rails   Opaque  foo             "dtan4\n"
```

### `k8sec unset`

Unset secrets

``` bash
$ k8sec unset [--kubeconfig KUBECONFIG] [--namespace NAMESPACE] NAME KEY1 KEY2

# Example
$ k8sec unset rails rails-env
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/dtan4/k8sec
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
