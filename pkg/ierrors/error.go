// Package ierrors TODO ADD PKG description
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

// New is the func similar to the standard library `errors.New` but it
// returns the inspr error structure, containing an error code and the
// capability of wrapping the message with extra context messages
func New(format string, values ...interface{}) *ierror { // MATCH /New.*unexported/
	return &ierror{
		err:  fmt.Errorf(format, values...),
		code: Unknown,
	}
}

// From is the func to create a ierror using as a base the an error interface
func From(err error) *ierror {
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
	return fmt.Sprintf("%v", ie.err.Error())
}

// Is has the purpose of establishing a way to utilize the standard library func
// known as "errors.Is(source,target)", by doing the overloading of the "Is"
// func it allows the comparison of an ierror with any other type of error.
//
// When using errors.Is(source Ierror, target error) it will return true if the
// `source` fully unwrapped is the same as `target` or if both of them have the
// same ErrCode.
func (ie *ierror) Is(err error) bool {
	// checks if is another type of error inside the error stack
	if errors.Is(ie.err, err) {
		return true
	}

	// converts target to ierror structure, if possible
	t, ok := err.(*ierror)
	if !ok {
		return false
	}

	return ie.code&t.code > 0
}

// Code is a function that tries to convert the error interface to an ierror
// structure and returns its code value
func Code(err error) ErrCode {
	ierr, ok := err.(*ierror)
	if !ok {
		// attaches the code unknown to the new ierror
		ierr = From(err)

	}
	return ierr.code
}

// TODO REVIEW wrap -> descrition and maybe the possibility of using multiple
// string values at once like:
//
// ierrors.Wrap(err, "msg1", "msg2", "msg3")

// Wrap is responsible for adding extra context to an error, this is done by
// stacking error messages that give the receiver of the error the path of the
// error occurrence
func Wrap(err error, format string, values ...interface{}) error {
	// returns nil if error doesn't exist
	if err == nil {
		return nil
	}

	// if not an ierror type, makes the conversion
	ierr, ok := err.(*ierror)
	if !ok {
		ierr = From(err)
	}

	msg := fmt.Sprintf(format, values...)

	// if not empty
	if msg != "" {
		ierr.err = fmt.Errorf("%v : %w", msg, ierr.err)
	}

	return ierr
}

// TODO REVIEW unwrap descrition

// Unwrap is a err function that is capable of handling both the standard golang
// error as well as the insp.error structure, it removes the last wrap done to
// the err stack and if that was the last error in the stack it will return nil.
func Unwrap(err error) error {
	ierr, ok := err.(*ierror)
	if !ok {
		return errors.Unwrap(err)
	}

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
	t := &parseStruct{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	ie.code = t.Code
	ie.err = stackError(t.Stack)
	return ie
}

// stackT.error converts the following structure of a error stack message
// into an actual stack of errors using the fmt.errorf
func stackError(stack string) error {
	var err error

	// reverses the stack to so they are inserted in the proper order
	messages := strings.Split(stack, ":")

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	for _, msg := range messages {
		m := strings.TrimSpace(msg)
		if err == nil {
			err = errors.New(m)
		} else {
			err = Wrap(err, m)
		}
	}

	return err
}

// The functions bellow are designed so the.codeCode of a ierror can only be
// modified in a way after being instantiated, that meaning that one should use
// the following:
//
// ierrors.New("my message").NotFound()
//
// This allows us to create functions to change the.codeCode state but it doesn't
// add exported functions to the pkg.

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
