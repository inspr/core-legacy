# ![logo](./docs/img/inspr_logo.png)Inspr

[![codecov](https://codecov.io/gh/inspr/inspr/branch/develop/graph/badge.svg?token=C8SPTHPXMG)](https://codecov.io/gh/inspr/inspr)

Inspr is an engine for running distributed applications, using multiple communication patterns such as pub sub and more, focused on type consistency and development simplicity.

- :muscle: Robust: built on top of golang, kubernetes and other state of the art technologies

- :sparkles: Distributed: created to allow complex and hierarchical distributed design patterns

- :cloud: Cloud native: lightweight and modular, built with interchangeable pieces

- :hammer_and_wrench: Versatile: can adapt to any cloud environment

## Installation

### Insprd

To install Insprd, add the helm chart repository using the following command:

`helm repo add inspr https://inspr-charts.storage.googleapis.com`

Install the latest version with

`helm install inspr_name inspr/insprd`

The default values file is included in the `build/helm` folder.

### CLI

To install the CLI get the latest release for your architecture from the [`releases`](https://github.com/inspr/inspr/-/releases) page and add it to your PATH.

## Documentation

You can check the documentation on the [Confluence page for Inspr](https://inspr.atlassian.net/wiki/spaces/INX/overview)

## License

> TODO
