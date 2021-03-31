// Package utils has a number of useful operations that are used
// in multiple places of the inspr packages, contains operations
// such as:
//
// 	- go-cmp: "comparators and evaluator for slices and maps"
// 	- string_slice: "set of operations of custom string slice"
//	- string_operations: "set of operations to check certain conditions a string variable"
package utils

// CheckEmptyChannel is a simple checkup if the unmarshaled
// channel name is empty of not, since it can be used for
// the error message is important to inform if it was given
// or not.
func CheckEmptyChannel(channel string) string {
	if channel == "" {
		return "<not given>"
	}
	return channel
}
