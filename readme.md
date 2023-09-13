# BRGEO

## Introduction

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

## Go version

The minimum required version for Go is `1.21.1` becasue of the inclusion of the `log/slog` package, used for structured logging within the application.

