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
		Err  error   `yaml:"_" json:"_"`
		Code ErrCode `yaml:"code"  json:"code"`
	}

	// parseStruct is used in the unmarshal and marshal of ierror struct
	parseStruct struct {
		Stack string  `yaml:"stack" json:"stack"`
		Code  ErrCode `yaml:"code"  json:"code"`
	}
)

// Error returns the ierror Message
func (err *ierror) Error() string {
	return err.fullMessage()
}

func (err *ierror) fullMessage() string {
	// TODO fmtState?

	return ""
}

// TODO REVIEW wrap descrition

// Wrap is responsible for adding extra context to an error, this is done by
// stacking error messages that give the receiver of the error the path of the
// error occurrence
func Wrap(err error, msg string) error {
	// returns nil if error doesn't exist
	if err == nil {
		return nil
	}

	// if not an ierror type, makes the conversion
	ierr, ok := err.(*ierror)
	if !ok {
		ierr = NewError().
			InnerError(err).
			Code(ExtenalPkgError).
			Build()
	}

	// checks if there is a '%w' wrapper in the msg, if not it will add it to
	// the end of the error message.
	// like ('my message: %w', err)
	if strings.Contains(msg, "%w") {
		return fmt.Errorf(msg, ierr)
	} else {
		msg += ": %w"
		return fmt.Errorf(msg, ierr)
	}
}

// Unwrap is a err function that is capable of handling both the standard golang
// error as well as the insprError structure, it removes the last wrap done to
// the err stack and if that was the last error in the stack it will return nil.
func Unwrap(err error) error {
	ierr, ok := err.(*ierror)
	if !ok {
		return errors.Unwrap(err)
	}

	// unwraps the insprError
	ierr.Err = errors.Unwrap(ierr.Err)

	// if there is no other error inside the inspr stack, returns nil
	if ierr.Err == nil {
		return nil
	}

	return ierr
}

// Is Compares errors
func (err *ierror) Is(target error) bool {

	// check if is another type of error inside the error stack
	if errors.Is(err.Err, target) {
		return true
	}

	// is it is an ierror it checks the code
	t, ok := target.(*ierror)
	if !ok {
		return false
	}

	return t.Code&err.Code > 0
}

// HasCode Compares error with error code
func (err *ierror) HasCode(code ErrCode) bool {
	return (code & err.Code) > 0
}

// MarshalJSON a struct function that allows for operations to be done
// before or after the json.Marshal procedure
func (err *ierror) MarshalJSON() ([]byte, error) {
	t := parseStruct{
		// TODO func to create the stack
		Stack: "unwrapped error",
		Code:  err.Code,
	}
	return json.Marshal(t)
}

// UnmarshalJSON a struct function that allows for operations to be done
// before or after the json.Unmarshal procedure
func (err *ierror) UnmarshalJSON(data []byte) error {
	t := &parseStruct{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	err.Code = t.Code
	err.Err = stackToError(t.Stack)
	return err
}

// stackToError converts the following structure of a error stack message
// into an actual stack of errors using the fmt.Errorf
func stackToError(stack string) error {
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
			err = fmt.Errorf("%v: %w", m, err)
		}
	}

	return err
}
