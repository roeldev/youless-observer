`youless-observer` Helm Charts
==============================

[{{ template "chart.versionBadge" . }}][release-url]
[![Artifact Hub][artifact-hub-img]][artifact-hub-url]

[release-url]: https://github.com/roeldev/youless-observer/releases/tag/v{{ template "chart.version" . }}

[artifact-hub-img]: https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/youless-observer

[artifact-hub-url]: https://artifacthub.io/packages/search?repo=youless-observer

{{ template "chart.description" . }}

## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```sh
helm repo add youless-observer https://roeldev.github.io/youless-observer
```

{{ template "chart.valuesSection" . }}

## License

Copyright © 2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in
the [LICENSE](https://github.com/roeldev/youless-observer/blob/main/LICENSE) file.
