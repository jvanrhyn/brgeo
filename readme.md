
# brgeo: Facade for KeyCDN Geolocation

[![tag](https://img.shields.io/github/tag/jvanrhyn/brgeo.svg)](https://github.com/jvanrhyn/brgeo/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/jvanrhyn/brgeo?status.svg)](https://pkg.go.dev/github.com/jvanrhyn/brgeo)
![Build Status](https://github.com/jvanrhyn/brgeo/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/jvanrhyn/brgeo)](https://goreportcard.com/report/github.com/jvanrhyn/brgeo)
[![Coverage](https://img.shields.io/codecov/c/github/jvanrhyn/brgeo)](https://codecov.io/gh/jvanrhyn/brgeo)
[![Contributors](https://img.shields.io/github/contributors/jvanrhyn/brgeo)](https://github.com/jvanrhyn/brgeo/graphs/contributors)
[![License](https://img.shields.io/github/license/jvanrhyn/brgeo)](./LICENSE)




## 🚀 Install

```sh
go get github.com/jvanrhyn/brgeo
```

**Compatibility**: go >= 1.21


## 💡 Usage

### Description

A simple facade library to wrap access to external IP Address geo-location information. Currently, it supports the following services:
KeyCDN: https://tools.keycdn.com/geo

Third party endpoint results is mapped to a well-known model, so you can easily switch between services.

```go
type LookupResponse struct {
	City        string `json:"city"`
	RegionName  string `json:"region"`
	CountryName string `json:"country"`
}
```

Query results are cached using the `github.com/patrickmn/go-cache` library. 


## 🤝 Contributing

- Ping me on mastodon [@jvanrhyn](https://mastodon.world/@jvanrhyn) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/jvanrhyn/brgeo)
- Fix [open issues](https://github.com/jvanrhyn/brgeo/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=jvanrhyn/brgeo)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/jvanrhyn?style=for-the-badge)](https://github.com/sponsors/jvanrhyn)

## 📝 License

Copyright © 2023 [Johan van Rhyn](https://github.com/jvanrhyn).

This project is [MIT](./LICENSE) licensed.
