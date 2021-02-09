package operator

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// toDeployment - formats data from a DApp to a
// string used in the kubernetes template
func toDeployment(app *meta.App) string {
	return fmt.Sprintf("%v", app.Meta.Name)
}

// TODO REVIEW THE FUNCTION ABOVE, should be unique
