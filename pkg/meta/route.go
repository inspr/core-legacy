package meta

import (
	"inspr.dev/inspr/pkg/utils"
)

// RouteConnection is the structure to the pod address and its endpoints
type RouteConnection struct {
	Address   string
	Endpoints utils.StringArray
}
