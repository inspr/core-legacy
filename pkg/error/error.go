package errors

import (
	"errors"
	"fmt"
)

// InsprError is an error that happened inside inspr
type InsprError struct {
	Value string         `json:"value,omitempty" xml:"value" avro:"value"  xml:"value"  avro:"value"`
	Err   error          `json:"-"`
	Code  InsprErrorCode `json:"code,omitempty" xml:"code" avro:"code"  xml:"code"  avro:"code"`
	Stack string         `json:"stack,omitempty"  xml:"stack"  avro:"stack"`
}

func (err *InsprError) Error() string {
	return err.Value
}

// NewNotFoundError creates a not found error for the given thing
func NewNotFoundError(name string, namespace string) *InsprError {
	return &InsprError{
		fmt.Sprintf("component %v in namespace %v not found.", name, namespace),
		nil,
		NotFound,
		"",
	}
}

// IsNotFound tests if an error is a not found error
func IsNotFound(err error) bool {
	if converted, ok := err.(*InsprError); ok {
		return converted.Code == NotFound
	}
	return false
}

// NewEncodingError returns an error message for ecoding errors
func NewEncodingError(msg string, innerError error) *InsprError {
	return &InsprError{
		msg,
		innerError,
		Encoding,
		innerError.Error(),
	}
}

// IsEncoding tests if an error is an encoding error
func IsEncoding(err error) (value bool) {
	if converted, ok := err.(*InsprError); ok {
		value = value || converted.Code == Encoding
	}
	return
}

// IsConfiguration tests if an error is a configuration error
func IsConfiguration(err error) (value bool) {
	if converted, ok := err.(*InsprError); ok {
		value = value || converted.Code == ChannelConfiguration
		value = value || converted.Code == PipelineConfiguration
		value = value || converted.Code == NodeConfiguration
	}
	return
}

// NewAlreadyExistsError creates a not found error for the given thing
func NewAlreadyExistsError(name string, namespace string) *InsprError {
	return &InsprError{
		fmt.Sprintf("component %v in namespace %v already exists.", name, namespace),
		nil,
		AlreadyExists,
		"",
	}
}

// IsAlreadyExists tests if an error is a not found error
func IsAlreadyExists(err error) bool {
	if converted, ok := err.(*InsprError); ok {
		return converted.Code == AlreadyExists
	}
	return false
}

// Unwrap unwrapps the error
func (err *InsprError) Unwrap() error {
	return errors.Unwrap(err.Err)
}

// Is Compares errors
func (err *InsprError) Is(target error) bool {
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}
	return t.Code == err.Code
}

// ToNative converts an error to native golang format
func (err *InsprError) ToNative() interface{} {
	return map[string]interface{}{
		"code": err.Code,
		"err": func() string {
			if err.Err == nil {
				return ""
			}
			return fmt.Sprint(err.Err)
		}(),
		"value": err.Value,
	}
}

// InsprErrorCode is error codes for inspr errors
type InsprErrorCode int32

// Error codes for inspr errors
const (
	NotFound InsprErrorCode = iota + 1
	Encoding
	AlreadyExists
	NodeConfiguration
	PipelineConfiguration
	ChannelConfiguration
	ServiceCreation
	DeploymentCreation
	ServiceDeletion
	Connection
)
