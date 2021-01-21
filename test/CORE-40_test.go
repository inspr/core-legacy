/*
DApps with managed lifecycles

As a dApp developer, I would like to create applications whose nodes (sub-apps) and channels can be created and destroyed at runtime so that my application only consume resources when needed.
*/
package inspr_example

import (
	. "github.com/ptcar2009/ginkgo"
)

var _ = Scenario("basic creation and deletion", func() {
	Given("a running inspr cluster", func() {
		When("there is a channel connected to a dApp", func() {
			And("a command is sent to delete this channel", func() {
				Then("that channel will not be deleted and the command sender will receive an error message stating that the channel is still connected to a dApp", func() {

				})
			})

			And("a command is sent to delete this channel ", func() {
				And("that command has a “force”  flag on it", func() {
					Then("channel and it’s connected dApps will be deleted", func() {

					})
				})
			})

		})

		When("there is a channel that is not connected to any dApp ", func() {
			And("a command is sent to delete this channel", func() {
				Then("that channel will be deleted", func() {

				})
			})
		})

		When("a command is sent to create a dApp ", func() {
			And("this dApp has connections stated that don’t match existing channels", func() {
				Then("the dApp will not be created and the command sender will receive an error message stating that the app has not been created.", func() {

				})
			})
		})

		When("a command is sent to create a dApp ", func() {
			And("the dApp definition conforms with all the requirements", func() {
				Then("that dApp will be created", func() {

				})
			})
		})

		When("a command is sent to create a channel", func() {
			Then("that channel will be created", func() {

			})
		})

		When("there is a running dApp", func() {
			And("a command is sent to delete this dApp", func() {
				Then("that dApp will be deleted", func() {

				})
			})
		})

	})
})

var _ = Scenario("Changes in the underlying structure", func() {
	Given("a running Inspr app", func() {
		When("a connection is changed by the controller in one of the sub-apps", func() {
			And("the connection is supported by the existing channels", func() {
				Then("the sub-app will have its connections changed", func() {

				})
			})

			And("the connection cannot be made due to non existing channels", func() {
				Then("the sub-app will NOT have its connections changed and an error message will be sent to the command sender", func() {

				})
			})

		})
	})
})
