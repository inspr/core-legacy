package yamls

import _ "embed" // needs to be imported for embed to work

//go:embed general.yaml
// PingPongYAML is the yaml of the primes example
var PingPongYAML string
