package ierrors

// InsprErrorCode is error codes for inspr errors
type InsprErrorCode int32

// Error codes for inspr errors
const (
	NotFound InsprErrorCode = iota + 1
	AlreadyExists
	InternalServer
	InvalidName
	InvalidChannel
	InvalidApp
	InvalidChannelType
	BadRequest
)
