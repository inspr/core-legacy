package ierrors

import "fmt"

// ErrBuilder is an Inspr Error Creator
type ErrBuilder struct {
	err *InsprError
}

// NewError is the start function to create a New Error
func NewError() *ErrBuilder {
	return &ErrBuilder{
		err: &InsprError{},
	}
}

// Message adds a message to the error
func (b *ErrBuilder) Message(format string, values ...interface{}) *ErrBuilder {
	b.err.Message = fmt.Sprintf(format, values...)
	b.err.Stack = b.err.Message
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

/*
From this point foward are the functions that uses the constants ierror values, these are merely functions that simplyfy the process of specifying the current
type of error.

Instead of using builder.SetErrorCode(ierrors.NotFound), one could simply use
the following builder.NotFound() which already does the process above.
*/

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

// Forbidden adds Forbidden code to Inspr Error
func (b *ErrBuilder) Forbidden() *ErrBuilder {
	b.err.Code = Forbidden
	return b
}

// Unauthorized adds Unauthorized code to Inspr Error
func (b *ErrBuilder) Unauthorized() *ErrBuilder {
	b.err.Code = Unauthorized
	return b
}
