// THIS IS THE MASTER
package main

import (
	"fmt"
)

type Order struct {
	ID     string           `diff:"id"`
	Items  []string         `diff:"items"`
	Orders map[string]Order `diff:"orders"`
}

// type Change struct {
// 	Type string      // The type of change detected; can be one of create, update or delete
// 	Path []string    // The path of the detected change; will contain any field name or array index that was part of the traversal
// 	From interface{} // The original value that was present in the "from" structure
// 	To   interface{} // The new value that was detected as a change in the "to" structure
// }

func main() {
	mapizao := map[string]map[string]bool{
		"channel": {},
		"app":     {},
		"ctype":   {},
	}
	mapizao["app"]["app1"] = true
	mapizao["app"]["app2"] = true

	printSet(mapizao["app"])
	// aux1 := Order{
	// 	ID:    "457",
	// 	Items: []string{"121", "341"},
	// }
	// mapizin := map[string]Order{
	// 	"channel": aux1,
	// 	"app":     aux1,
	// 	"ctype":   aux1,
	// }
	// fmt.Println(mapizao["app"]["app3"])
	// fmt.Println(mapizao)
	// fmt.Println("Len App: ", len(mapizao["app"]))
	// fmt.Println("Len Chann: ", len(mapizao["channel"]))
	// fmt.Println(mapizin["appe"].ID == "")
	// aux1 := Order{
	// 	ID:    "4567",
	// 	Items: []string{"121", "341"},
	// }

	// aux2 := Order{
	// 	ID:    "457",
	// 	Items: []string{"121", "341"},
	// }

	// a := Order{
	// 	ID:     "1234",
	// 	Items:  []string{"um", "dois", "tres"},
	// 	Orders: map[string]Order{"aux1": aux1},
	// }

	// b := Order{
	// 	ID:     "1234",
	// 	Items:  []string{"um", "dois", "tres"},
	// 	Orders: map[string]Order{"aux1": aux2},
	// }
	// a := map[string]Order{"aux1": aux2}
	// b := map[string]Order{"aux2": aux1}
	// changelog, err := diff.Diff(a, b)

	// if err != nil {
	// 	fmt.Println(err, changelog)
	// }

	// aStr := []string{"um", "cinco", "quatro", "sete"}
	// bStr := []string{"asdf", "fdas", "fds", "asdf"}
	// fullStr := append(aStr, bStr...)
	// fmt.Println(fullStr)
	// for _, elem := range changelog {
	// 	to := fmt.Sprintf("%v", elem.To)
	// 	fmt.Println("Type: ", elem.Type)
	// 	fmt.Println("Path: ", elem.Path)
	// 	fmt.Println("From: ", elem.From)
	// 	fmt.Println("To: ", to)
	// }
	// fmt.Println(reflect.TypeOf(changelog[0].To))
	// for _, elem := range changelog[0].Path {
	// 	fmt.Println(elem)
	// }
	// fmt.Println(len(returnsNilSlice()))
}

func returnsNilSlice() []string {
	return nil
}

func printSet(set map[string]bool) {
	for index := range set {
		fmt.Println("Index: ", index)
		// fmt.Println("Elem: ", elem)
	}
}
