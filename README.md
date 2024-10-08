# testman

😎 `go test` wrapper for advanced testing workflows in Go

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/comnoco/testman/v2)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/moul/testman/blob/master/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/moul/testman.svg)](https://github.com/moul/testman/releases)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/testman.svg)](https://microbadger.com/images/moul/testman)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)

[![Go](https://github.com/moul/testman/workflows/Go/badge.svg)](https://github.com/moul/testman/actions?query=workflow%3AGo)
[![Release](https://github.com/moul/testman/workflows/Release/badge.svg)](https://github.com/moul/testman/actions?query=workflow%3ARelease)
[![PR](https://github.com/moul/testman/workflows/PR/badge.svg)](https://github.com/moul/testman/actions?query=workflow%3APR)
[![GolangCI](https://golangci.com/badges/github.com/moul/testman.svg)](https://golangci.com/r/github.com/moul/testman)
[![codecov](https://codecov.io/gh/moul/testman/branch/master/graph/badge.svg)](https://codecov.io/gh/moul/testman)
[![Go Report Card](https://goreportcard.com/badge/github.com/comnoco/testman/v2)](https://goreportcard.com/report/github.com/comnoco/testman/v2)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/testman/badge)](https://www.codefactor.io/repository/github/moul/testman)


## Usage

*testman -h*

[embedmd]:# (.tmp/root-usage.txt)
```txt
USAGE
  testman <subcommand> [flags]

SUBCOMMANDS
  test  advanced go test workflows
  list  list available tests
```

*testman test -h*

[embedmd]:# (.tmp/test-usage.txt)
```txt
USAGE
  testman test [flags] [packages]

EXAMPLES
   testman test ./...
   testman test -v ./...
   testman test -run ^TestUnstable -timeout=300s -retry=50 ./...
   testman test -run ^TestBroken -test.timeout=30s -retry=10 --continue-on-error ./...
   testman test -test.timeout=10s -test.v -test.count=2 -test.race

FLAGS
  -continue-on-error false  continue on error (but still fails at the end)
  -retry 0                  fail after N retries
  -run ^(Test|Example)      regex to filter out tests and examples
  -test.count 1             `go test -count=VAL`
  -test.race false          `go test -race`
  -test.timeout 0s          `go test -timeout=VAL`
  -test.v false             `go test -v`
  -timeout 0s               program max duration
  -v false                  verbose
```

*testman list -h*

[embedmd]:# (.tmp/list-usage.txt)
```txt
USAGE
  testman list [packages]

EXAMPLES
   testman list ./...
   testman list -v ./...
   testman list -run ^TestStable ./...

FLAGS
  -run ^(Test|Example)  regex to filter out tests and examples
  -v false              verbose
```

## Install

### Using go

```console
$ go get -u github.com/comnoco/testman/v2
```

### Releases

See https://github.com/comnoco/testman/releases

## Contribute

![Contribute <3](https://raw.githubusercontent.com/moul/moul/master/contribute.gif)

I really welcome contributions. Your input is the most precious material. I'm well aware of that and I thank you in advance. Everyone is encouraged to look at what they can do on their own scale; no effort is too small.

Everything on contribution is sum up here: [CONTRIBUTING.md](./CONTRIBUTING.md)

### Contributors ✨

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://manfred.life"><img src="https://avatars1.githubusercontent.com/u/94029?v=4" width="100px;" alt=""/><br /><sub><b>Manfred Touron</b></sub></a><br /><a href="#maintenance-moul" title="Maintenance">🚧</a> <a href="https://github.com/moul/testman/commits?author=moul" title="Documentation">📖</a> <a href="https://github.com/moul/testman/commits?author=moul" title="Tests">⚠️</a> <a href="https://github.com/moul/testman/commits?author=moul" title="Code">💻</a></td>
    <td align="center"><a href="https://manfred.life/moul-bot"><img src="https://avatars1.githubusercontent.com/u/41326314?v=4" width="100px;" alt=""/><br /><sub><b>moul-bot</b></sub></a><br /><a href="#maintenance-moul-bot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

### Stargazers over time

[![Stargazers over time](https://starchart.cc/moul/testman.svg)](https://starchart.cc/moul/testman)

## License

© 2020-2021 [Manfred Touron](https://manfred.life)

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT) ([`LICENSE-MIT`](LICENSE-MIT)), at your option. See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`
