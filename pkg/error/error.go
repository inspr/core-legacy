package errors

import "fmt"

// InsprError is an error that happened inside inspr
type InsprError struct {
	Message string
	Err     error
	Code    InsprErrorCode
}

// Error returns the InsprError Message
func (err *InsprError) Error() string {
	return err.Message
}

// NewCustomError provides a method to create a custom error, given the error code and the error message
func NewCustomError(errCode InsprErrorCode, errMsg string) *InsprError {
	return &InsprError{
		errMsg,
		nil,
		errCode,
	}
}

// NewNotFoundError creates a new Not Found Inspr Error
func NewNotFoundError(name string, err error) *InsprError {
	return &InsprError{
		fmt.Sprintf("Component %v not found.", name),
		err,
		NotFound,
	}
}

// NewAlreadyExistsError creates a new Already Exists Inspr Error
func NewAlreadyExistsError(name string, err error) *InsprError {
	return &InsprError{
		fmt.Sprintf("Component %v already exists.", name),
		err,
		AlreadyExists,
	}
}

// NewInternalServerError creates a new Internal Server Inspr Error
func NewInternalServerError(err error) *InsprError {
	return &InsprError{
		fmt.Sprintf("There was a internal server error."),
		err,
		InternalServer,
	}
}

// NewInvalidNameError creates a new Invalid Name Inspr Error
func NewInvalidNameError(name string, err error) *InsprError {
	return &InsprError{
		fmt.Sprintf("The name '%v' is invalid.", name),
		err,
		InvalidName,
	}
}

// NewInvalidChannelError creates a new Invalid Channel Inspr Error
func NewInvalidChannelError() *InsprError {
	return &InsprError{
		fmt.Sprintf("The channel is invalid."),
		nil,
		InvalidChannel,
	}
}

// NewInvalidAppError creates a new Invalid App Inspr Error
func NewInvalidAppError() *InsprError {
	return &InsprError{
		fmt.Sprintf("The app is invalid."),
		nil,
		InvalidApp,
	}
}

// NewInvalidChannelTypeError creates a new Invalid ChannelType Inspr Error
func NewInvalidChannelTypeError() *InsprError {
	return &InsprError{
		fmt.Sprintf("The ChannelType is invalid."),
		nil,
		InvalidChannelType,
	}
}

// Is Compares errors
func (err *InsprError) Is(target error) bool {
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}
	return t.Code == err.Code
}

// HasCode Compares error with error code
func (err *InsprError) HasCode(code InsprErrorCode) bool {
	return code == err.Code
}
