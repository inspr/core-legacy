package ierrors

// InsprErrorCode is error codes for inspr errors
type InsprErrorCode int32

// Error codes for inspr errors
const (
	NotFound InsprErrorCode = 1 << iota
	AlreadyExists
	InternalServer
	InvalidName
	InvalidApp
	InvalidChannel
	InvalidChannelType
	BadRequest
	Unauthorized
	Forbidden
)
