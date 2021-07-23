package ierrors

// Code is error codes for inspr errors
type ErrCode uint32

// Error codes for inspr errors
const (
	NotFound ErrCode = 1 << iota
	AlreadyExists
	InternalServer
	InvalidName
	InvalidApp
	InvalidChannel
	InvalidType
	InvalidFile
	InvalidArgs
	BadRequest
	InvalidToken
	ExpiredToken
	Unauthorized
	Forbidden
	ExtenalPkgError
)
