<p align="center">
    <img 
        srcset="./docs/img/logo.png, ./docs/img/logo@2x.png, ./docs/img/logo@3x.png "
        src="./docs/img/logo.png" 
        width="128" 
        title="Inspr Logo">
</p>

<p align="center">
    <a href="https://godoc.org/github.com/inspr/inspr">
        <img src="https://godoc.org/github.com/inspr/inspr?status.svg">
    </a>
    <a href="https://goreportcard.com/badge/github.com/inspr/inspr">
        <img src="https://goreportcard.com/badge/github.com/inspr/inspr">
    </a>
    <a href="https://codecov.io/gh/inspr/inspr">
        <img src="https://codecov.io/gh/inspr/inspr/branch/develop/graph/badge.svg?token=C8SPTHPXMG">
    </a>
</p>

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

To install the CLI get the latest release for your architecture from the [`releases`](https://github.com/inspr/inspr/releases) page and add it to your PATH.

## Documentation

You can check the documentation on the [Confluence page for Inspr](https://inspr.atlassian.net/wiki/spaces/INX/overview)

