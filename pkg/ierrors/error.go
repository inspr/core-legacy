package ierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// InsprError is an error that happened inside inspr
type InsprError struct {
	Message string         `yaml:"message"  json:"message"`
	Err     error          `yaml:"_" json:"_"`
	Stack   string         `yaml:"stack" json:"stack"`
	Code    InsprErrorCode `yaml:"code"  json:"code"`
}

// Error returns the InsprError Message
func (ierror *InsprError) Error() string {
	return ierror.Message
}

// Is Compares errors
func (ierror *InsprError) Is(target error) bool {

	// check if is another type of error inside the error stack
	if errors.Is(ierror.Err, target) {
		return true
	}

	// is it is an InsprError it checks the code
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}

	return t.Code&ierror.Code > 0
}

// HasCode Compares error with error code
func (ierror *InsprError) HasCode(code InsprErrorCode) bool {
	return code == ierror.Code
}

// Wrap adds a message to the ierror stack
func (ierror *InsprError) Wrap(message string) {
	if ierror.Err == nil {
		ierror.Err = errors.New(message)
	} else {
		ierror.Err = fmt.Errorf("%v: %w", message, ierror.Err)
	}
	ierror.Message = ierror.Err.Error()
	ierror.Stack = ierror.Err.Error()
}

// Wrapf adds the format with the values given to the ierror stack
func (ierror *InsprError) Wrapf(format string, values ...interface{}) {
	message := fmt.Sprintf(format, values...)
	ierror.Wrap(message)
}

// MarshalJSON a struct function that allows for operations to be done
// before or after the json.Marshal procedure
func (ierror *InsprError) MarshalJSON() ([]byte, error) {
	return json.Marshal(*ierror)
}

// UnmarshalJSON a struct function that allows for operations to be done
// before or after the json.Unmarshal procedure
func (ierror *InsprError) UnmarshalJSON(data []byte) error {
	t := struct {
		Message string         `yaml:"message"  json:"message"`
		Stack   string         `yaml:"stack" json:"stack"`
		Code    InsprErrorCode `yaml:"code"  json:"code"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	// copies it to the insprErr
	ierror.Message = t.Message
	ierror.Stack = t.Stack
	ierror.Code = t.Code
	ierror.StackToError()
	return nil
}

// StackToError converts the following structure of a error stack message
// into an actual stack of errors using the fmt.Errorf
func (ierror *InsprError) StackToError() {
	messages := strings.Split(ierror.Stack, ":")

	// reverses the stack to so they are inserted in the proper order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	for _, msg := range messages {
		m := strings.TrimSpace(msg)
		if ierror.Err == nil {
			ierror.Err = errors.New(m)
		} else {
			ierror.Err = fmt.Errorf("%v: %w", m, ierror.Err)
		}
	}
}
