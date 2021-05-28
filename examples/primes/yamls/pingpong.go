package yamls

import _ "embed" // needs to be imported for embed to work

//go:embed general.yaml

// PrimesYAML is the yaml of the primes example
var PrimesYAML string
