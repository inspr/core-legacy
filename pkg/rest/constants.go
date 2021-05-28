package rest

import "net/http"

const (
	// HeaderScopeKey is used to obtain or to change the scope in the http
	// request header, usually used in the handlers of the api package and the
	// client of the controller package.
	HeaderScopeKey = "Scope"
)

// operationTranslator receives a http.Method and proceeds to translate into a
// string that represents the operation that will be done in the insprd
var operationTranslator = map[string]string{
	http.MethodGet:    "get",
	http.MethodPost:   "create",
	http.MethodPut:    "update",
	http.MethodDelete: "delete",
}

// routeTranslator receives a route that was processed from a http.URL and then
// proceeds to translate into a string that represents its abbreviation used for
// the const permission in the auth pkg.
var routeTranslator = map[string]string{
	"apps":     "dapp",
	"channels": "channel",
	"types":    "type",
	"auth":     "token",
}
