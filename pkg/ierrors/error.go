// package ierrors
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

func (ie *ierror) Is(err error) bool {
	// check if is another type of error inside the error stack
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

// New is the function to create a New.error
func New(format string, values ...interface{}) *ierror {
	return &ierror{
		err:  fmt.Errorf(format, values...),
		code: Unknown,
	}
}

// From is the function to create a ierror using as a base the an error interface
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
func (err *ierror) Error() string {
	return fmt.Sprintf("%v", err.err.Error())
}

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
func (err *ierror) MarshalJSON() ([]byte, error) {

	// there is no way of setting the inner error as nil using the exported funcs,
	// one would have to set it inside the ierrors pkg.
	if err.err == nil {
		return []byte{}, New("unexpected err, ierror inner error field got set to nil").ExternalErr()
	}

	t := parseStruct{
		Stack: err.err.Error(),
		Code:  err.code,
	}
	return json.Marshal(t)
}

// UnmarshalJSON a struct function that allows for operations to be done
// before or after the json.Unmarshal procedure
func (err *ierror) UnmarshalJSON(data []byte) error {
	t := &parseStruct{}

	json.Unmarshal(data, &t)

	err.code = t.Code
	err.err = stackError(t.Stack)
	return err
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
func (e *ierror) NotFound() *ierror {
	e.code = NotFound
	return e
}

// AlreadyExists adds Already Exists code to Inspr Error
func (e *ierror) AlreadyExists() *ierror {
	e.code = AlreadyExists
	return e
}

// BadRequest adds Bad Request code to Inspr Error
func (e *ierror) BadRequest() *ierror {
	e.code = BadRequest
	return e
}

// InternalServer adds Internal Server code to Inspr Error
func (e *ierror) InternalServer() *ierror {
	e.code = InternalServer
	return e
}

// InvalidName adds Invalid Name code to Inspr Error
func (e *ierror) InvalidName() *ierror {
	e.code = InvalidName
	return e
}

// InvalidApp adds Invalid App code to Inspr Error
func (e *ierror) InvalidApp() *ierror {
	e.code = InvalidApp
	return e
}

// InvalidChannel adds Invalid Channel code to Inspr Error
func (e *ierror) InvalidChannel() *ierror {
	e.code = InvalidChannel
	return e
}

// InvalidType adds Invalid Type code to Inspr Error
func (e *ierror) InvalidType() *ierror {
	e.code = InvalidType
	return e
}

// InvalidFile adds Invalid Args code to Inspr Error
func (e *ierror) InvalidFile() *ierror {
	e.code = InvalidFile
	return e
}

// InvalidToken adds Invalid Token code to Inspr Error
func (e *ierror) InvalidToken() *ierror {
	e.code = InvalidToken
	return e
}

// InvalidArgs adds Invalid Args code to Inspr Error
func (e *ierror) InvalidArgs() *ierror {
	e.code = InvalidArgs
	return e
}

// Forbidden adds Forbidden code to Inspr Error
func (e *ierror) Forbidden() *ierror {
	e.code = Forbidden
	return e
}

// Unauthorized adds Unauthorized code to Inspr Error
func (e *ierror) Unauthorized() *ierror {
	e.code = Unauthorized
	return e
}

// ExternalErr adds ExternalPkgError code to Inspr Error
func (e *ierror) ExternalErr() *ierror {
	e.code = ExternalPkg
	return e
}
