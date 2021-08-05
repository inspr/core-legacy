// Package ierrors provides an easy way to wrap context to the error of the
// standard library, as well as support for multiple errors.
//
// Declaring an error
//
// This is the most usual way of using the package, there are two options: the
// first consists of using ierrors.New(""), specifying the message of the error,
// and the ierrors.MultiError structure whichs supports multiple errors.
//
// - ierror: is the structure created by the ierrors.New() and consists of a
// message and an errorCode, the second parameter can be added when
// instantiating an ierror after the New func call. An example of its usage would
// be `ierrors.New("cli subcommand %v error", subcommandName).InvalidArgs()`.
//
// - MultiError: structure which is composed of an errorCode and a slice of
// error interfaces, it contains the methods Add(error) and Empty(). Which can
// be used to handle error requests that are comming from a set of goroutines.
//
// Some functionalities that are contained within this package are:
//
// - Wrap/Unwrap functions: handles the addition and removal of context to an
// error, meaning that if desired the developer can provide extra information to
// the error he can simply use `ierrors.Wrap(err, "on operation X")`.
//
// - Marshal/Unmarshal: the ierror structure declared by this package have
// custom functions that are called when using json.Marshal or json.Unmarshal,
// meaning that when creating a error it can be parsed into json and yaml
// format, allowing the transfer of error information via http.
//
// - New: as described previously is used to handle the creation of errors, one
// difference from the standard package is that this func accepts strings and
// other errors as a parameter. Meaning that one could create a inspr error from
// another error of an external pkg.
//
// Some examples of the usage of the `New` func are:
//	- ierrors.New("my_custom_err")
//	- ierrors.New("error on the URL %v and route %v", myURL, routeName)
//	- ierrors.New(io.EOF) // error defined in the standard library io pkg
//
// - ierror.is(): an internal func that is called when using the standard
// library errors.Is function. This allows us to not lose the original error and
// compare it at any point of our code. For example when using the error below:
//
// errWithWrapping := ierrors.Wrap(io.EOF, "func X", "file Y")
//
// We can still check at any point if the base of the error is an io.EOF, that
// is done by simply calling errors.Is(errWithWrapping, io.EOF). No matter how
// many wrappers you attach to the original error the contents can still be
// compared to another error.
//
package ierrors // import "inspr.dev/inspr/pkg/ierrors"
