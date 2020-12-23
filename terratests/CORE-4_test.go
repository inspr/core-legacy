/*
Second Template of task, as example

h1. This is another task

As a developer, I want to have a template task so that I know what to do on Jira
*/
package inspr_example

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe(" 3: is the context of the sprint", func() {
	Context(" a new sprint", func() {
		When(" that sprint starts"+
			"&& that sprint has tasks", func() {
			It(" the tasks layouts will match this", func() {

			})
		})
	})
})

var _ = Describe(" 4: is the context of the task", func() {
	Context(" a task", func() {
		When(" that task contains Functional Requirements"+
			"&& those requirements are in given when then format", func() {
			It(" the requirements will be translated to gingko tests", func() {

			})
		})
	})
})
