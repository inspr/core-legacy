/*
Functional dApps

As a developer, I want to transform my applications into k8s’s deployments, so that it can run on the cloud.
*/
package inspr_example

import (
	. "github.com/ptcar2009/ginkgo"
)

var _ = Scenario("give a dApp deployment without any channel attached", func() {
	Given("that I have a running cluster on Inspr", func() {
		When("I deploy my application on that cluster *and* my application is a set of other apps", func() {
			Then("all the application’s subApps will be deployed", func() {

			})
		})

		When("I deploy my application on that cluster *and* my application is a single node", func() {
			Then("that application will be deployed as a k8s deployment", func() {

			})
		})

	})
})

var _ = Scenario("channels inside a dApp", func() {
	Given("a running cluster on Inspr", func() {
		When("I deploy my application with a channel attached to it", func() {
			Then("the application will be deployed *and* the channels will be created in the specified message broker", func() {

			})
		})
	})
})

var _ = Scenario("compliance", func() {
	Given("a running Inspr cluster", func() {
		When("a command is sent to it", func() {
			Then("the commands should work as described in [*CORE-40*|https://inspr.atlassian.net/browse/CORE-40]", func() {

			})
		})
	})
})
