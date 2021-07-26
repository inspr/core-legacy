package ierrors

// The functions bellow are designed so the.CodeCode of a ierror can only be
// modified in a way after being instantiated, that meaning that one should use
// the following:
//
// ierrors.New("my message").NotFound()
//
// This allows us to create functions to change the.CodeCode state but it doesn't
// add exported functions to the pkg.

// NotFound adds Not Found code to Inspr Error
func (e *ierror) NotFound() *ierror {
	e.Code = NotFound
	return e
}

// AlreadyExists adds Already Exists code to Inspr Error
func (e *ierror) AlreadyExists() *ierror {
	e.Code = AlreadyExists
	return e
}

// BadRequest adds Bad Request code to Inspr Error
func (e *ierror) BadRequest() *ierror {
	e.Code = BadRequest
	return e
}

// InternalServer adds Internal Server code to Inspr Error
func (e *ierror) InternalServer() *ierror {
	e.Code = InternalServer
	return e
}

// InvalidName adds Invalid Name code to Inspr Error
func (e *ierror) InvalidName() *ierror {
	e.Code = InvalidName
	return e
}

// InvalidApp adds Invalid App code to Inspr Error
func (e *ierror) InvalidApp() *ierror {
	e.Code = InvalidApp
	return e
}

// InvalidChannel adds Invalid Channel code to Inspr Error
func (e *ierror) InvalidChannel() *ierror {
	e.Code = InvalidChannel
	return e
}

// InvalidType adds Invalid Type code to Inspr Error
func (e *ierror) InvalidType() *ierror {
	e.Code = InvalidType
	return e
}

// InvalidFile adds Invalid Args code to Inspr Error
func (e *ierror) InvalidFile() *ierror {
	e.Code = InvalidFile
	return e
}

// InvalidArgs adds Invalid Args code to Inspr Error
func (e *ierror) InvalidArgs() *ierror {
	e.Code = InvalidArgs
	return e
}

// Forbidden adds Forbidden code to Inspr Error
func (e *ierror) Forbidden() *ierror {
	e.Code = Forbidden
	return e
}

// Unauthorized adds Unauthorized code to Inspr Error
func (e *ierror) Unauthorized() *ierror {
	e.Code = Unauthorized
	return e
}

// ExternalErr adds ExternalPkgError code to Inspr Error
func (e *ierror) ExternalErr() *ierror {
	e.Code = ExtenalPkg
	return e
}
