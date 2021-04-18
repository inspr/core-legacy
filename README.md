<br>
<br>
<div align="center">
    <img 
        srcset="https://storage.googleapis.com/inspr/logo.png 1x, https://storage.googleapis.com/inspr/logo%402x.png 2x, https://storage.googleapis.com/inspr/logo%403x.png 3x"
        src="https://storage.googleapis.com/inspr/logo%403x.png" 
        width="128" 
        title="Inspr Logo"></img>
</div>

<h2 align="center">A Cloud Native Application Mesh.</h2>
<p align="center">Describe, connect, secure, test and deploy distributed applications with less effort.</p>

<p align="center">
    <a href="https://godoc.org/github.com/inspr/inspr">
        <img src="https://godoc.org/github.com/inspr/inspr?status.svg">
    </a>
    <a href="https://goreportcard.com/badge/github.com/inspr/inspr">
        <img src="https://goreportcard.com/badge/github.com/inspr/inspr">
    </a>
    <a href="https://pkg.go.dev/github.com/inspr/inspr">
        <img src="https://img.shields.io/static/v1?label=godoc&message=reference&color=blue">
    </a>
    <a href="https://codecov.io/gh/inspr/inspr">
        <img src="https://codecov.io/gh/inspr/inspr/branch/develop/graph/badge.svg?token=C8SPTHPXMG&">
    </a>
    <a href="https://github.com/inspr/inspr/blob/develop/LICENSE.md">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg">
    </a>
    <a href="http://makeapullrequest.com">
        <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg">
    </a>
    <a href="https://discord.gg/tmp2564a54">
        <img src="https://img.shields.io/discord/449569975888248832.svg?label=&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2">
    </a>
</p>

<!-- <p align="center">
    <a href="https://stars.medv.io/inspr/inspr">
        <img src="https://stars.medv.io/inspr/inspr.svg">
    </a>
</p> -->


---

<br>
Inspr is an engine for running distributed applications, using multiple communication patterns such as pub sub and more, focused on type consistency and development simplicity.

- :muscle: Robust: built on top of golang, kubernetes and other state of the art technologies
- :sparkles: Distributed: created to allow complex and hierarchical distributed design patterns
- :cloud: Cloud native: lightweight and modular, built with interchangeable pieces
- :hammer_and_wrench: Versatile: can adapt to any cloud environment

## Table of Contents

- [Getting Started](#getting-started)
- [Documentation](#documentation)
- [Architecture](#architecture)
- [Contributing](#-contributing)
- [License](#-license)
- [Contact](#-contact)

## 🎉 Getting Started

### Cluster

The recommended way to install inspr in a kubernetes cluster is by using helm.

the first step is add the helm chart repository to the cluster:

```bash
helm repo add inspr https://inspr-charts.storage.googleapis.com
```

Then install inspr with the command:

```bash
helm install inspr_name inspr/insprd
```

replacing inspr_name by the desired inspr cluster name.

Aditionaly you can check the default values file for the helm chart. 
They are included in the `build/helm` folder and can be edited for further refinement of the properties.

### CLI

To install the CLI get the latest release for your architecture from the [`releases`](https://github.com/inspr/inspr/releases) page and add it to your PATH.
TODO: Add link to the install script for the CLI.


## Documentation

The full documentation for Inspr can be found on our [website](https://inspr.com/docs).
The source for the Inspr documentation and website is hosted on a separate repo, @inspr/docs.

## Architecture
- TBD

## 👏 Contributing

The main purpose of this repository is to continue evolving Inspr core. We want to make contributing to this project as easy and transparent as possible, and we are grateful to the community for contributing bug fixes and improvements.
To contribute, please check out our [guide](CONTRIBUTING.md).

## 📄 License

- Inspr is licenced under MIT [license](LICENSE.md).
- Inspr documentation is Creative Commons licensed, as found in the [license-docs](LICENSE-docs) file.

## 📣 Contact

To contact us, please join our [Discord community](https://discord.gg/tmp2564a54).
Aditionally we are at [Reddit](https://reddit.com/r/inspr)

