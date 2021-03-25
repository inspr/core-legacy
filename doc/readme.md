# Getting Started

This file is directed for those who are running the inspr plataform for the first time, it is recommended that you take your time through each step and read carefully all the documentation linked below.

## Helm

To be able to use the Inspr daemon there is the requirement that you have a cluster available to use. There is also a possibility to use [minikube](https://minikube.sigs.k8s.io/docs/start/) to test it locally.

For the step by step installation of all the necessary software in the cluster please follow this [tutorial](./helm_installation.md).

## CLI

After preparing your cluster it's necessary to install the Inspr CLI, this can be done though this [tutorial](./cli_install.md).

## Example

Now that all the necessary installation steps are taken, it's time to test the capabilities of the Inspr plataform.

Any of the examples in the `examples` folder are available for testing but it is recommended that in your first time the `primes` example is used. For that you can either clone the repository or download just the primes folder directly.

The steps to building your application are the following:
- Open the terminal
- Enter in the directory of the primes folder
- Run the command `make` to build the application docker image and push it to the Inspr repository.
- Run `inspr apply -f yamls/general.yaml` to install the primes example in your cluster, for a futher understanting of what this command do you can always use `inspr apply -h`.
- Use k9s to check the state of your cluster and see if there are pods running with the name `print`, `generator` and `filter`.
  - if so check the logs of each of them to see what is being written.
  - Generator should be logging random number
  - Filter is responsible for taking the primes created by the Generator
  - Printer is responsible for printing all the primes received by the Filter
