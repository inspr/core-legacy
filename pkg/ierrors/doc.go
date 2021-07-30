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
// be `ierrors.New("cli subcommand X error").InvalidArgs()`.
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
// - From: it receives an external error interace and converts it to an ierror
// structure while mantaining the content. For example
// ierrors.From(sql.ErrNoRows) creates an ierror with the base error of the sql
// package, meaning that is possible to add futher context to the sql message.
// By doing this it is still possible to use the errors.Is(err_1, err_2) from
// the standard library to compare different types of errors even after adding
// the contexts.
package ierrors // import "inspr.dev/inspr/pkg/ierrors"
