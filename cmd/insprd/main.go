// THIS IS THE MASTER
package main

import "gitlab.inspr.dev/inspr/core/cmd/insprd/api"

type person struct {
	name string
	age  int
}

func main() {
	api.Run()
}
