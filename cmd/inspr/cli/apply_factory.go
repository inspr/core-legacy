package cli

import (
	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

/*
RunMethod defines the method that will run for the
component
*/
type RunMethod func(*cobra.Command, []string)

/*
ApplyFactory holds a dictionary that maps all the pairs
(kind, apiVersion - encapsulated in the Component Type)
to a cobra run method
*/
type ApplyFactory struct {
	applyDict map[meta.Component]RunMethod
}

var applyFactory *ApplyFactory

/*
GetFactory returns the ApllyFactory singleton.
If it doesn't exist, create one
*/
func GetFactory() *ApplyFactory {
	if applyFactory == nil {
		applyFactory = &ApplyFactory{
			applyDict: make(map[meta.Component]RunMethod),
		}
	}
	return applyFactory
}

/*
GetRunMethod returns the runMethod registered for the
given component. If the component is not found in the
dictionary, it returns a ierror
*/
func (af *ApplyFactory) GetRunMethod(component meta.Component) (RunMethod, error) {
	if method, ok := af.applyDict[component]; ok {
		return method, nil
	}
	return nil, ierrors.NewError().
		InvalidName().
		Message("Component not subscribed in the ApplyFactory dictionary").
		Build()
}

/*
Subscribe adds to the apply factory dictonary the
given component with the value equals to the given
runMethod
*/
func (af *ApplyFactory) Subscribe(component meta.Component, method RunMethod) {
	af.applyDict[component] = method
}
