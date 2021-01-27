// THIS IS THE MASTER
package main

import (
	"fmt"
	"strings"
)

func main() {
	var strtest string = "tralala.trololo"
	if strings.Contains(strtest, ".") {
		fmt.Println("TEM SIM")
	}
}
