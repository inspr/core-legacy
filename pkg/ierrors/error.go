package ierrors

import "fmt"

// InsprError is an error that happened inside inspr
type InsprError struct {
	Message string         `yaml:"message"  json:"message"`
	Err     error          `yaml:"_" json:"_"`
	Code    InsprErrorCode `yaml:"code"  json:"code"`
}

// ErrBuilder is an Inspr Error Creator
type ErrBuilder struct {
	err *InsprError
}

// Error returns the InsprError Message
func (err *InsprError) Error() string {
	return err.Message
}

// NewError is the start function to create a New Error
func NewError() *ErrBuilder {
	return &ErrBuilder{
		err: &InsprError{},
	}
}

// NotFound adds Not Found code to Inspr Error
func (b *ErrBuilder) NotFound() *ErrBuilder {
	b.err.Code = NotFound
	return b
}

// AlreadyExists adds Already Exists code to Inspr Error
func (b *ErrBuilder) AlreadyExists() *ErrBuilder {
	b.err.Code = AlreadyExists
	return b
}

// BadRequest adds Bad Request code to Inspr Error
func (b *ErrBuilder) BadRequest() *ErrBuilder {
	b.err.Code = BadRequest
	return b
}

// InternalServer adds Internal Server code to Inspr Error
func (b *ErrBuilder) InternalServer() *ErrBuilder {
	b.err.Code = InternalServer
	return b
}

// InvalidName adds Invalid Name code to Inspr Error
func (b *ErrBuilder) InvalidName() *ErrBuilder {
	b.err.Code = InvalidName
	return b
}

// InvalidApp adds Invalid App code to Inspr Error
func (b *ErrBuilder) InvalidApp() *ErrBuilder {
	b.err.Code = InvalidApp
	return b
}

// InvalidChannel adds Invalid Channel code to Inspr Error
func (b *ErrBuilder) InvalidChannel() *ErrBuilder {
	b.err.Code = InvalidChannel
	return b
}

// InvalidChannelType adds Invalid Channel Type code to Inspr Error
func (b *ErrBuilder) InvalidChannelType() *ErrBuilder {
	b.err.Code = InvalidChannelType
	return b
}

// Message adds a message to the error
func (b *ErrBuilder) Message(format string, values ...interface{}) *ErrBuilder {
	b.err.Message = fmt.Sprintf(format, values...)
	return b
}

// InnerError adds a inner error to the error
func (b *ErrBuilder) InnerError(err error) *ErrBuilder {
	b.err.Err = err
	return b
}

// Build returns the created Inspr Error
func (b *ErrBuilder) Build() *InsprError {
	return b.err
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

func HasCode(target error, code InsprErrorCode) bool {
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}
	return t.Code&code > 0
}
