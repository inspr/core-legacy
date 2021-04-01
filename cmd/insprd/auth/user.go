package auth

const ( // User roles
	standard = 0
	admin    = 1
)

//User is an object containing all info about an Inspr user
type User struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Role     int      `json:"role"`
	Scopes   []string `json:"scopes"`
	Token    string   `json:"token"`
	UID      string   `json:"uid"`
}

//Builder interface for building an User
type Builder interface {
	SetName(name string) Builder
	SetPassword(pwd string) Builder
	AsAdmin() Builder
	SetScope(scope ...string) Builder
}

//internal builder structure
type builder struct {
	usr User
}

//NewUser instanciates a standard user on the internal builder
func NewUser() Builder {
	return &builder{
		usr: User{
			Scopes: make([]string, 0),
			Role:   standard,
		},
	}
}

//SetName names the user being built
func (bd *builder) SetName(name string) Builder {
	bd.usr.Name = name
	return bd
}

//SetPassword configures the user's password
func (bd *builder) SetPassword(pwd string) Builder {
	bd.usr.Password = pwd
	return bd
}

//SetScopes defines which scopes the user is allowed access to
func (bd *builder) SetScope(scope ...string) Builder {
	bd.usr.Scopes = scope
	return bd
}

//AsAdmin configures the user as an Admin
func (bd *builder) AsAdmin() Builder {
	bd.usr.Role = admin
	return bd
}
