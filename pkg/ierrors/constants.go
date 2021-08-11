package ierrors

// ErrCode is error codes for inspr errors
type ErrCode uint32

// Error codes for inspr errors
const (
	Unknown ErrCode = 1 << iota
	NotFound
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
	ExternalPkg
)
