package ierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/disiqueira/gotree"
)

// InsprError is an error that happened inside inspr
type InsprError struct {
	Message string         `yaml:"message"  json:"message"`
	Err     error          `yaml:"_" json:"_"`
	Stack   string         `yaml:"stack" json:"stack"`
	Code    InsprErrorCode `yaml:"code"  json:"code"`
}

// Error returns the InsprError Message
func (err *InsprError) Error() string {
	return err.Message
}

// Is Compares errors
func (err *InsprError) Is(target error) bool {
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}
	return t.Code&err.Code > 0
}

// HasCode Compares error with error code
func (err *InsprError) HasCode(code InsprErrorCode) bool {
	return code == err.Code
}

// TODO TESTS
// Wrap adds a message to the ierror stack
func (ierror *InsprError) Wrap(message string) {
	ierror.Err = fmt.Errorf("%v: %w", message, ierror.Err)
	ierror.Stack = ierror.Err.Error()
}

// TODO TESTS
// Wrapf adds the format with the values given to the ierror stack
func (ierror *InsprError) Wrapf(format string, values ...interface{}) {
	message := fmt.Sprintf(format, values...)
	ierror.Err = fmt.Errorf("%v: %w", message, ierror.Err)
}

// TODO TESTS
// MarshalJSON a struct function that allows for operations to be done
// before or after the json.Marshal procedure
func (ierror *InsprError) MarshalJSON() ([]byte, error) {
	return json.Marshal(ierror)
}

// UnmarshalJSON a struct function that allows for operations to be done
// before or after the json.Unmarshal procedure
func (ierror *InsprError) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &ierror); err != nil {
		return err
	}
	ierror.StackToError()
	return nil
}

// StackToError converts the following structure of a error stack message
// into an actual stack of errors using the fmt.Errorf
func (ierror *InsprError) StackToError() {
	messages := strings.Split(ierror.Stack, ":")
	for _, msg := range messages {
		m := strings.TrimSpace(msg)
		if ierror.Err == nil {
			ierror.Err = errors.New(m)
		} else {
			ierror.Err = fmt.Errorf("%v: %w", m, ierror.Err)
		}
	}
}

// FormatedError the main focus of this function is to allow for a more
// readable error information, so it will be used mainly in debug sessions
func (ierror *InsprError) FormatedError() {
	tree := gotree.New("ErrorTree")
	messages := strings.Split(ierror.Err.Error(), ":")
	for _, msg := range messages {
		m := strings.TrimSpace(msg)
		tree.Add(m)
	}
}
