package ierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type (
	// ierror is an error that happened inside inspr
	ierror struct {
		err  error
		code ErrCode
	}

	// parseStruct is used in the unmarshal and marshal of ierror struct
	parseStruct struct {
		Stack string  `yaml:"stack" json:"stack"`
		Code  ErrCode `yaml:"code"  json:"code"`
	}
)

const (
	// prefixMessage is used by the stackToErr to remove the error message
	// before processing the stack of errors
	prefixMessage = "error :"

	// separator is the content that will be between the error and each
	// additional context inserted by a Wrap message
	separator = ":"
)

// New is the function used to create an ierror, it accepts an interface{} since
// an ierror can be created from one of the following:
//
// - A string value
// - An error interface
//
// If none of the above are provided as a argument the function will return nil.
//
// One useful way of handling errors from an external pkg is to use this
// function to translate it to an ierror, allowing the customization of the
// ErrCode.
//
// Some examples of the usage of the `New` func are:
//	- ierrors.New("my_custom_err")
//	- ierrors.New("error on the URL %v and route %v", myURL, routeName)
//	- ierrors.New(io.EOF) // error defined in the standard library io pkg
//
func New(param interface{}, a ...interface{}) *ierror {
	var ie *ierror

	switch v := param.(type) {
	case error:
		ie = from(v)
	case string:
		ie = newIerror(v, a...)
	default:
		ie = nil
	}

	return ie
}

// newIerror is the func similar to the standard library `errors.New` but it
// returns the inspr error structure, containing an error code and the
// capability of wrapping the message with extra context messages
func newIerror(format string, values ...interface{}) *ierror {
	var msg string
	if len(values) > 0 {
		msg = fmt.Sprintf(format, values...)
	} else {
		msg = format
	}

	return &ierror{
		err:  errors.New(msg),
		code: Unknown,
	}
}

// from is the func to create an ierror structure using as a base the an error
// interface
func from(err error) *ierror {
	ierr, ok := err.(*ierror)
	if !ok {
		ierr = &ierror{
			err:  err,
			code: Unknown,
		}
	}
	return ierr
}

// Error returns the ierror Message
func (ie *ierror) Error() string {
	return fmt.Sprintf("%v %v", prefixMessage, ie.err)
}

// FormatError is a simple function with the intention of handling the default
// ierror Error() format into something more presentable.
//
// Meaning that the error
//
// "file X : func Y : <base_error> " will be converted into:
//
// file X
//		func Y
//		<base_error>
func FormatError(err error) string {
	ie := from(err)
	stack := strings.TrimPrefix(ie.err.Error(), prefixMessage)

	// reverses the stack to so they are inserted in the proper order
	messages := strings.Split(stack, ":")

	var message string
	for i, msg := range messages {
		m := strings.TrimSpace(msg)
		if i == 0 {
			message += fmt.Sprintf("%v %v\n", prefixMessage, m)
		} else {
			message += fmt.Sprintf("\t%v\n", m)
		}
	}
	return message
}

// Is has the purpose of establishing a way to utilize the standard library func
// known as "errors.Is(source,target)", by doing the overloading of the "Is"
// func it allows the comparison of an ierror with any other type of error.
//
// When using errors.Is(source Ierror, target error) it will return true if the
// `source` fully unwrapped is the same as `target` or if both of them have the
// same ErrCode.
//
// An observation is that by default ierrors with the Unknown error code will
// not be considered equal to another ierror with that code.
func (ie *ierror) Is(err error) bool {
	// checks if is another type of error inside the error stack
	if errors.Is(ie.err, err) {
		return true
	}

	// when using an ierror with unkown error in the comparison the result will
	// always be false
	if ie.code == Unknown {
		return false
	}

	// converts target to ierror structure, if possible
	t := from(err)

	return ie.code&t.code > 0
}

// Code is a function that tries to convert the error interface to an ierror
// structure and returns its code value
func Code(err error) ErrCode {
	return from(err).code
}

// Wrap is responsible for adding extra context to an error, this is done by
// stacking error messages that give the receiver of the error the path of the
// error occurrence
func Wrap(err error, msgs ...string) error {
	// returns nil if error doesn't exist
	if err == nil {
		return nil
	}

	// if not an ierror type, makes the conversion
	ierr := from(err)

	for _, msg := range msgs {
		if msg != "" {
			ierr.err = fmt.Errorf("%s %s %w", msg, separator, ierr.err)
		}
	}

	return ierr
}

// Unwrap is a err function that removes the last wrap done to the err stack
// and if that was the last error in the stack it will return nil. It is
// capable of handling both the standard golang error as well as the insprError
// structure.
func Unwrap(err error) error {
	ierr := from(err)

	// unwraps the insp.error
	ierr.err = errors.Unwrap(ierr.err)

	// if there is no other error inside the inspr stack, returns nil
	if ierr.err == nil {
		return nil
	}

	return ierr
}

// MarshalJSON a struct function that allows for operations to be done
// before or after the json.Marshal procedure
func (ie *ierror) MarshalJSON() ([]byte, error) {

	// there is no way of setting the inner error as nil using the exported funcs,
	// one would have to set it inside the ierrors pkg.
	if ie.err == nil {
		return []byte{},
			New("unexpected err, ierror inner error field got set to nil").ExternalErr()
	}

	t := parseStruct{
		Stack: ie.err.Error(),
		Code:  ie.code,
	}
	return json.Marshal(t)
}

// UnmarshalJSON a struct function that allows for operations to be done
// before or after the json.Unmarshal procedure
func (ie *ierror) UnmarshalJSON(data []byte) error {
	t := &parseStruct{
		Code: Unknown,
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	parsedIerr := stackToError(t.Stack, t.Code)
	ie.err = parsedIerr.err
	ie.code = parsedIerr.code
	return nil
}

// stackToError converts a stack message and a code into an ierror strucutre
func stackToError(stack string, code ErrCode) *ierror {
	var ie *ierror

	stack = strings.TrimPrefix(stack, prefixMessage)

	// reverses the stack to so they are inserted in the proper order
	messages := strings.Split(stack, separator)

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	for _, msg := range messages {
		msg = strings.TrimSpace(msg)

		if ie == nil {
			ie = New(msg)
			ie.code = code
		} else {
			ie.err = fmt.Errorf("%v %s %w", msg, separator, ie.err)
		}
	}

	return ie
}

// The functions bellow are designed so the.codeCode of a ierror can only be
// modified in a way after being instantiated, that meaning that one should use
// the following:
//
// ierrors.New("my message").NotFound()
//
// This allows us to create functions to change the ierr.Code state but it
// doesn't add exported variables to the error structure.

// NotFound adds Not Found code to Inspr Error
func (ie *ierror) NotFound() *ierror {
	ie.code = NotFound
	return ie
}

// AlreadyExists adds Already Exists code to Inspr Error
func (ie *ierror) AlreadyExists() *ierror {
	ie.code = AlreadyExists
	return ie
}

// BadRequest adds Bad Request code to Inspr Error
func (ie *ierror) BadRequest() *ierror {
	ie.code = BadRequest
	return ie
}

// InternalServer adds Internal Server code to Inspr Error
func (ie *ierror) InternalServer() *ierror {
	ie.code = InternalServer
	return ie
}

// InvalidName adds Invalid Name code to Inspr Error
func (ie *ierror) InvalidName() *ierror {
	ie.code = InvalidName
	return ie
}

// InvalidApp adds Invalid App code to Inspr Error
func (ie *ierror) InvalidApp() *ierror {
	ie.code = InvalidApp
	return ie
}

// InvalidChannel adds Invalid Channel code to Inspr Error
func (ie *ierror) InvalidChannel() *ierror {
	ie.code = InvalidChannel
	return ie
}

// InvalidType adds Invalid Type code to Inspr Error
func (ie *ierror) InvalidType() *ierror {
	ie.code = InvalidType
	return ie
}

// InvalidFile adds Invalid Args code to Inspr Error
func (ie *ierror) InvalidFile() *ierror {
	ie.code = InvalidFile
	return ie
}

// InvalidToken adds Invalid Token code to Inspr Error
func (ie *ierror) InvalidToken() *ierror {
	ie.code = InvalidToken
	return ie
}

// InvalidArgs adds Invalid Args code to Inspr Error
func (ie *ierror) InvalidArgs() *ierror {
	ie.code = InvalidArgs
	return ie
}

// Forbidden adds Forbidden code to Inspr Error
func (ie *ierror) Forbidden() *ierror {
	ie.code = Forbidden
	return ie
}

// Unauthorized adds Unauthorized code to Inspr Error
func (ie *ierror) Unauthorized() *ierror {
	ie.code = Unauthorized
	return ie
}

// ExternalErr adds ExternalPkgError code to Inspr Error
func (ie *ierror) ExternalErr() *ierror {
	ie.code = ExternalPkg
	return ie
}
