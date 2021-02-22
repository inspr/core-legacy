// Package applymethods is the package with the functions
// to be inserted into the factory method
package applymethods

/* INTEGRATION
Given the current state of the tasks some things will change when integrating everything,
this is a list of things to do  related to this package

- remove the RunMethod and use the one defined in the method Factory
- Establisha a prerun in the apply cmd where the factory.Subscribe is done for this method
- how to determine
*/

// RunMethod defines the method that will run for the component
type RunMethod func([]byte) error

// ApplyChannel is of the type RunMethod, it calls the pkg/controller/client functions.
var ApplyChannel RunMethod = func([]byte) error {

	// unmarshal into a channel

	// creates or updates it

	return nil
}
