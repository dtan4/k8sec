# k8sec

CLI tool to manage [Kubernetes Secrets](http://kubernetes.io/docs/user-guide/secrets/) easily.

## Usage

### `k8sec list`

List secrets

``` bash
$ k8sec list [--namespace NAMESPACE] [--kubeconfig KUBECONFIG] [NAME]

# Example
$ k8sec list rails
NAME    TYPE    KEY             VALUE
rails   Opaque  database-url    "postgres://example.com:5432/dbname"
```

### `k8sec set`

Set secrets

``` bash
$ k8sec set [--namespace NAMESPACE] [--kubeconfig KUBECONFIG] NAME KEY1=VALUE1 KEY2=VALUE2

# Example
$ k8sec set rails RAILS_ENV=production
```

### `k8sec unset`

Unset secrets

``` bash
$ k8sec unset [--namespace NAMESPACE] [--kubeconfig KUBECONFIG] NAME KEY1 KEY2

# Example
$ k8sec unset rails RAILS_ENV
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
