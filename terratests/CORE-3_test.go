/*
Template task for gleison

h1. This is a task

As a developer, I want to have a template task so that I know what to do on Jira
*/
package inspr_example

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe(" 1: is the context of the sprint", func() {
	Context(" a new sprint", func() {
		When(" that sprint starts"+
			"&& that sprint has tasks", func() {
			It(" the tasks layouts will match this", func() {

			})
		})
	})
})

var _ = Describe(" 2: is the context of the task", func() {
	Context(" a task", func() {
		When(" that task contains Functional Requirements"+
			"&& those requirements are in given when then format", func() {
			It(" the requirements will be translated to gingko testsdd", func() {

			})
		})
	})
})
