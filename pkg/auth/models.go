package auth

// ResfreshDO is the body type expected by UID provider to refresh a payload
type ResfreshDO struct {
	RefreshToken []byte `json:"refreshtoken"`
}

// JwtDO is a data output type for master's token handler
type JwtDO struct {
	Token []byte `json:"token"`
}

// InitDO  structure for initialization requests
type InitDO struct {
	Payload
	Key string
}

//Payload is information caried by a Inspr acceess token
type Payload struct {
	UID string `json:"uid"`
	// Permissions is a map where key is the Scope and values are permissions
	Permissions map[string][]string `json:"permissions"`
	Refresh     []byte              `json:"refresh"`
	RefreshURL  string              `json:"refreshurl"`
}

/* TODO
	- Invert the permissions scheme. map will be indexed by the permission and the value should be a list of scopes.
		- Allows for the implementation of general permissions
		- Less compromising for other authentication structures and methods
		- Is acctually a lot more efficent size wise
	- Check validate behaviour for general functions
	- Check validate for init
	- Check validade for token creation
	- Check validade for brokers
	- Move logic of validation from middleware. To Auth?
TODONE */

// All Permissions possible values
const (
	CreateDapp    string = "create:dapp"
	CreateChannel string = "create:channel"
	CreateType    string = "create:type"
	CreateAlias   string = "create:alias"
	CreateBroker  string = "create:broker"

	GetDapp    string = "get:dapp"
	GetChannel string = "get:channel"
	GetType    string = "get:type"
	GetAlias   string = "get:alias"
	GetBroker  string = "get:broker"

	UpdateDapp    string = "update:dapp"
	UpdateChannel string = "update:channel"
	UpdateType    string = "update:type"
	UpdateAlias   string = "update:alias"

	DeleteDapp    string = "delete:dapp"
	DeleteChannel string = "delete:channel"
	DeleteType    string = "delete:type"
	DeleteAlias   string = "delete:alias"

	CreateToken string = "create:token"
)

// AdminPermissions defines all the permissions that the admin user have
// when the cluster is initialized
var AdminPermissions = map[string][]string{
	"": {
		CreateDapp,
		CreateChannel,
		CreateType,
		CreateAlias,
		CreateBroker,

		GetDapp,
		GetChannel,
		GetType,
		GetAlias,
		GetBroker,

		UpdateDapp,
		UpdateChannel,
		UpdateType,
		UpdateAlias,

		DeleteDapp,
		DeleteChannel,
		DeleteType,
		DeleteAlias,

		CreateToken,
	},
}
