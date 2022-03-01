# [v0.8.0](https://github.com/dtan4/k8sec/releases/tag/v0.8.0) (2022-03-02)

## Features

- Customize User-Agent with k8sec version ([#153](https://github.com/dtan4/k8sec/pull/153)) (thanks @bendrucker)

## Others

- Use Go 1.17
- Update dependencies including some security fixes
- Stop providing official Docker image
- Drop support under Kubernetes 1.18 (see also: [Kubernetes version skew policy](https://kubernetes.io/releases/version-skew-policy/))

# [v0.7.0](https://github.com/dtan4/k8sec/releases/tag/v0.7.0) (2020-03-13)

- Use Go 1.14
- Update dependencies to the latest one (e.g. client-go v0.17.3)

# [v0.6.0](https://github.com/dtan4/k8sec/releases/tag/v0.6.0) (2018-08-05)

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
