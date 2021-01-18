/*
Modular core interfaces

As a developer of Inspr I would like my code to be composable, so that the next modules to be developed have reduced integration cost and a faster time to market.
*/
package inspr_example

import (
	. "github.com/onsi/ginkgo"
)

var _ = Context("Scenario 1: our development team built a new feature for Inspr", func() {
	Given("*Given* a new feature was developed", func() {
		And("*and* passed all the bug test", func() {
			When("*When* the team wants to deploy that feature to a small group for user test", func() {
				Then("*Then* the team just launch it without any trouble or bug of the features that already existed on Inspr.", func() {
				})
			})
		})
	})
})
