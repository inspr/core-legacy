package yamls

import _ "embed" // needs to be imported for embed to work

// PrimesYAML is the yaml of the primes example
//go:embed general.yaml
var PrimesYAML string
