# [v0.7.0](https://github.com/dtan4/k8sec/releases/tag/v0.7.0) (2019-04-09)

## Features

- Sort secrets on dump subcommand [#1](https://github.com/wantedly/k8sec/pull/1)

## Others

- Fork repository [#2](https://github.com/wantedly/k8sec/pull/2)

# [v0.6.0](https://github.com/dtan4/k8sec/releases/tag/v0.5.1) (2018-08-05)

## Features

- Enable external auth providers [#32](https://github.com/dtan4/k8stail/pull/32)

## Fixed

- Exclude namespace from kubeclient [#28](https://github.com/dtan4/k8sec/pull/28)

## Others

- Use Go 1.10.3 on Travis CI [#34](https://github.com/dtan4/k8stail/pull/34)
- Upgrade to client-go 8.0.0 [#31](https://github.com/dtan4/k8stail/pull/31)

# [v0.5.1](https://github.com/dtan4/k8sec/releases/tag/v0.5.1) (2017-08-25)

## Fixed

- Sort `k8sec list` command output [#25](https://github.com/dtan4/k8sec/pull/25) (thanks @unblee)

# [v0.5.0](https://github.com/dtan4/k8sec/releases/tag/v0.5.0) (2017-08-15)

## Features

- Add `-n` flag as an alias of `--namespace` [#23](https://github.com/dtan4/k8sec/pull/23)

## Fixed

- Modify application name in version command [#22](https://github.com/dtan4/k8sec/pull/22)

# [v0.4.1](https://github.com/dtan4/k8sec/releases/tag/v0.4.1) (2017-04-12)

## Fixed

- Use correct namespace [#20](https://github.com/dtan4/k8sec/pull/20)

# [v0.4.0](https://github.com/dtan4/k8sec/releases/tag/v0.4.0) (2017-04-12)

## Features

- Select context / Use namespace set in kubecfg [#18](https://github.com/dtan4/k8sec/pull/18)

# [v0.3.1](https://github.com/dtan4/k8sec/releases/tag/v0.3.1) (2017-01-10)

## Fixed

- Update command description [#15](https://github.com/dtan4/k8sec/pull/15)
- Suppress usage and error printing at error [#13](https://github.com/dtan4/k8sec/pull/13)
- Check the length of key-value array [#12](https://github.com/dtan4/k8sec/pull/12)
- Create new secret if it does not exist [#11](https://github.com/dtan4/k8sec/pull/11)

# [v0.3.0](https://github.com/dtan4/k8sec/releases/tag/v0.3.0) (2016-12-30)

## Backward imcompatible changes

- Rename `k8sec save` command as `k8sec dump` [#9](https://github.com/dtan4/k8sec/pull/9)
  - No behavior changes.

# [v0.2.0](https://github.com/dtan4/k8sec/releases/tag/v0.2.0) (2016-12-07)

Drop Kubernetes <= 1.2 support

# [v0.1.1](https://github.com/dtan4/k8sec/releases/tag/v0.1.1) (2016-10-21)

Initial stable binary release.

# [v0.1.0](https://github.com/dtan4/k8sec/releases/tag/v0.1.0) (2016-07-19)

Initial release.
