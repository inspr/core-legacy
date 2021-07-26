package ierrors

// The functions bellow are designed so the.codeCode of a ierror can only be
// modified in a way after being instantiated, that meaning that one should use
// the following:
//
// ierrors.New("my message").NotFound()
//
// This allows us to create functions to change the.codeCode state but it doesn't
// add exported functions to the pkg.

// NotFound adds Not Found code to Inspr Error
func (e *ierror) NotFound() *ierror {
	e.code = NotFound
	return e
}

// AlreadyExists adds Already Exists code to Inspr Error
func (e *ierror) AlreadyExists() *ierror {
	e.code = AlreadyExists
	return e
}

// BadRequest adds Bad Request code to Inspr Error
func (e *ierror) BadRequest() *ierror {
	e.code = BadRequest
	return e
}

// InternalServer adds Internal Server code to Inspr Error
func (e *ierror) InternalServer() *ierror {
	e.code = InternalServer
	return e
}

// InvalidName adds Invalid Name code to Inspr Error
func (e *ierror) InvalidName() *ierror {
	e.code = InvalidName
	return e
}

// InvalidApp adds Invalid App code to Inspr Error
func (e *ierror) InvalidApp() *ierror {
	e.code = InvalidApp
	return e
}

// InvalidChannel adds Invalid Channel code to Inspr Error
func (e *ierror) InvalidChannel() *ierror {
	e.code = InvalidChannel
	return e
}

// InvalidType adds Invalid Type code to Inspr Error
func (e *ierror) InvalidType() *ierror {
	e.code = InvalidType
	return e
}

// InvalidFile adds Invalid Args code to Inspr Error
func (e *ierror) InvalidFile() *ierror {
	e.code = InvalidFile
	return e
}

// InvalidArgs adds Invalid Args code to Inspr Error
func (e *ierror) InvalidArgs() *ierror {
	e.code = InvalidArgs
	return e
}

// Forbidden adds Forbidden code to Inspr Error
func (e *ierror) Forbidden() *ierror {
	e.code = Forbidden
	return e
}

// Unauthorized adds Unauthorized code to Inspr Error
func (e *ierror) Unauthorized() *ierror {
	e.code = Unauthorized
	return e
}

// ExternalErr adds ExternalPkgError code to Inspr Error
func (e *ierror) ExternalErr() *ierror {
	e.code = ExternalPkg
	return e
}
