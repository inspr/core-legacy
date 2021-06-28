package yamls

import _ "embed" // needs to be imported for embed to work

//go:embed general.yaml
var PrimesYAML string // PrimesYAML is the yaml of the primes example
