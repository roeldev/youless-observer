youless-observer
================

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Image size][image-size-img]][image-size-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/roeldev/youless-logger.svg?label=latest

[latest-release-url]: https://github.com/roeldev/youless-logger/releases

[build-status-img]: https://github.com/roeldev/youless-logger/actions/workflows/test.yml/badge.svg

[build-status-url]: https://github.com/roeldev/youless-logger/actions/workflows/test.yml

[image-size-img]: https://img.shields.io/docker/image-size/roeldev/youless-observer

[image-size-url]: https://hub.docker.com/repository/docker/roeldev/youless-observer/tags

[report-img]: https://goreportcard.com/badge/github.com/roeldev/youless-logger

[report-url]: https://goreportcard.com/report/github.com/roeldev/youless-logger

[doc-img]: https://godoc.org/github.com/roeldev/youless-logger?status.svg

[doc-url]: https://pkg.go.dev/github.com/roeldev/youless-logger

Service `youless-observer` is a real-time utility usage observer for the _YouLess energy meter_ device. It collects
metrics data and can send it to any OTEL compatible backend, notify MQTT broker(s), or stream it via gRPC.

Key features are:

- export metrics via OTLP
- scrape metrics via Prometheus
- send metrics via MQTT (work in progress)
- stream metrics via gRPC (work in progress)

> Need to store the metrics data for a longer period of time? Have a look
> at [youless-logger](https://github.com/roeldev/youless-logger).

## Installation

The recommended way to run `youless-observer` is by using the available [Docker image](https://hub.docker.com/repository/docker/roeldev/youless-observer).
A [Helm chart](charts/youless-observer/README.md) is also available to easily deploy the service inside a Kubernetes cluster.

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=roeldev" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
