`youless-observer` Helm Charts
==============================

[![Version: 0.1.0-rc2](https://img.shields.io/badge/Version-0.1.0--rc2-informational?style=flat) ][release-url]
[![Artifact Hub][artifact-hub-img]][artifact-hub-url]

[release-url]: https://github.com/roeldev/youless-observer/releases/tag/v0.1.0-rc2

[artifact-hub-img]: https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/youless-observer

[artifact-hub-url]: https://artifacthub.io/packages/search?repo=youless-observer

Real-time usage observer service for YouLess energy meter

## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```sh
helm repo add youless-observer https://roeldev.github.io/youless-observer
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| image.repository | string | `"roeldev/youless-observer"` |  |
| image.tag | string | `"0.1.0-rc2"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.pullSecrets | list | `[]` |  |
| service.enabled | bool | `true` |  |
| service.type | string | `"LoadBalancer"` |  |
| service.port | int | `80` |  |
| service.annotations | object | `{}` |  |
| networkPolicy.create | bool | `true` |  |
| observer.debug | bool | `false` |  |
| observer.log.level | string | `"warn"` |  |
| observer.log.timestamp | bool | `false` |  |
| observer.log.accessLog | bool | `false` |  |
| observer.server.port | int | `2512` |  |
| observer.server.tls.enabled | bool | `false` |  |
| observer.server.tls.secretName | string | `""` |  |
| observer.server.tls.caSecretName | string | `""` |  |
| observer.server.tls.certFile | string | `""` |  |
| observer.server.tls.keyFile | string | `""` |  |
| observer.server.tls.caCertFile | string | `""` |  |
| observer.server.tls.verifyClient | bool | `true` |  |
| observer.server.tls.insecureSkipVerify | bool | `false` |  |
| observer.otel.serviceName | string | `""` |  |
| observer.otel.resourceAttributes | string | `""` |  |
| observer.otel.metric.exportInterval | int | `10000` |  |
| observer.otel.metric.exportTimeout | int | `10000` |  |
| observer.otel.traces.enabled | bool | `true` |  |
| observer.otel.traces.sampler | string | `"parentbased_traceidratio"` |  |
| observer.otel.traces.samplerArg | float | `0.5` |  |
| observer.otel.exporter.otlp.enabled | bool | `true` |  |
| observer.otel.exporter.otlp.endpoint | string | `""` |  |
| observer.otel.exporter.otlp.headers | string | `""` |  |
| observer.otel.exporter.otlp.protocol | string | `"grpc"` |  |
| observer.otel.exporter.otlp.timeout | int | `5000` |  |
| observer.otel.exporter.otlp.tls.enabled | bool | `false` |  |
| observer.otel.exporter.otlp.tls.secretName | string | `""` |  |
| observer.otel.exporter.otlp.tls.caSecretName | string | `""` |  |
| observer.otel.exporter.otlp.tls.certFile | string | `""` |  |
| observer.otel.exporter.otlp.tls.keyFile | string | `""` |  |
| observer.otel.exporter.otlp.tls.caCertFile | string | `""` |  |
| observer.otel.exporter.otlp.tls.insecureSkipVerify | bool | `false` |  |
| observer.prometheus.enabled | bool | `false` |  |
| observer.prometheus.endpoint | string | `"/metrics"` |  |
| observer.youless.configMapName | string | `""` |  |
| observer.youless.secretName | string | `""` |  |
| observer.youless.url | string | `"http://youless"` |  |
| observer.youless.name | string | `"YouLess"` |  |
| observer.youless.timeout | string | `"5s"` |  |
| observer.youless.password | string | `""` |  |
| observer.excludePower | bool | `false` |  |
| observer.excludeS0 | bool | `false` |  |
| observer.excludeGas | bool | `false` |  |
| observer.excludeWater | bool | `false` |  |
| observer.singlePhase | bool | `true` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| securityContext | object | `{}` |  |
| resources.limits.cpu | string | `"50m"` |  |
| resources.limits.memory | string | `"32Mi"` |  |
| resources.requests.cpu | string | `"50m"` |  |
| resources.memory | string | `"32Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |

## License

Copyright © 2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in
the [LICENSE](https://github.com/roeldev/youless-observer/blob/main/LICENSE) file.
