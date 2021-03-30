package auth

type User struct {
	Name   string   `json:"name"`
	Pwd    string   `json:"pwd"`
	Role   int      `json:"role"`
	Scopes []string `json:"scopes"`
	Token  string
	UID    string
}

type Builder interface {
	SetName(name string) Builder
	SetPassword(pwd string) Builder
	AsAdmin() Builder
	SetScope(scope ...string) Builder
}

type builder struct {
	usr User
}

func NewUser() Builder {
	return &builder{
		usr: User{
			Scopes: make([]string, 0),
		},
	}
}

func (bd *builder) SetName(name string) Builder {
	bd.usr.Name = name
	return bd
}

func (bd *builder) SetPassword(pwd string) Builder {
	bd.usr.Pwd = pwd
	return bd
}

func (bd *builder) SetScope(scope ...string) Builder {
	bd.usr.Scopes = scope
	return bd
}

func (bd *builder) AsAdmin() Builder {
	bd.usr.Role = 1
	return bd
}
