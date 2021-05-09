# Changelog
### #48 - Feature: adds init command to the cli

- features:
    - creates an init command for configuration initialization
    - adds config flag to specify configuration file
    - changes flow so that if the config is invalid, asks for initialization
- refactors:
    - refactors configuration initialization to allow different files
    - changes command builder to now use variadic flag argument instead of slice
- fix:
    - silences main command to remove duplicate errors


### #41 - Feature: Handling invalid requests in the middleware authorization
- features:
    - rest/middleware.go -> using the `permissions` field in the payload, created in the issue
      CORE-336, for the scope and permissions check.
    - reading the scope information from the http.request header, instead of the
      body
    - changing the controller/client so is more functional in nature, functions
      return a copy of the struct with the new value.
    - rest/request -> changed so the scope is inserted in the header of the
      request, the structs that were the body of the request were changed so
they no longer have a `scope` field.

- fixes:
    - fixed tests in the controller/client and cli/handlers that were broken due
      to some struct changes.

- tests:
    - new tests for the middleware, since there is a new logic step due to the
      permission introduction.


### #37 - Feature: reactive sidecar <!-- This is the title -->
- fixes:
	- lowers the timeout on kafka writer's flush
	- changes some loggers to be production and lower test cluttering
	- some type fixes
- features:
	- sidecar server now has a goroutine for constant message polling. This routine reads from all subscribed channels and sends requests to the main dApp so that the dapp handles those requests
	- sidecar's read message now has a maximum number of retries
	- sidecar and dApp now communicate via tcp instead of unix sockets. This is helpful to add easier integration with libs for many languages, including requests for js, Flask for python, etc. 
	- inspr client now acts as a wrapper on a http server for reading messages
	- write message now takes the routes as channel parameters
	- adds fields for port definition on the dApp metadata
- tests:
	- changes kafka tests to handle the new channel based approach

### #39 - Issue CORE-336 | Alter the token to contain the permissions... <!-- This is the title -->
- features:
	- changed Payload (`pkg/auth/models/payload.go`) field 'Role' (`int`) to 'Permission' (`map[string][]string`), and removed field 'Scope'
        - 'Permission' is now a map of Scope to list of Permissions
    - changed User (`cmd/uid_provider/client/interface.go`) field 'Role' (`int`) to 'Permission' (`map[string][]string`), and removed field 'Scope'
        - 'Permission' is now a map of Scope to list of Permissions
    - Adapted all methods which use the structures above:
        - `cmd/uid_provider/client/client.go`
        - `cmd/uid_provider/inprov/cmd/create_user.go`
        - `pkg/api/handlers/token_handler.go`
        - `pkg/auth/mocks/auth_mock.go`
        - `pkg/controller/client/auth.go`
        - `pkg/rest/middleware.go`
        - `pkg/api/handlers/token_handler.go`
- fixes:
	- Inspr CLI token flag default value is now being properly set
	- removed duplicated Payload struct that was being defined in `pkg/api/auth/payload.go`
- refactor:
    - grouped files under `pkg/auth/models` into `models.go` in `pkg/auth`, and deleted the folder (this was done so the structures have the package prefix `auth` when being imported)
        - Adapted all methods which use the structures defined in `models` so they continue to work as intended
- tests:
	- changes in tests of methods which were impacted by Payload and User modification
- misc:
    - removed `init()` func that wasn't working in `cmd/insprd/operators/kafka/nodes/converter.go` (@ptcar2009 fixed it PR #37)
    - created example yaml file for uid provider user's structure definition (`cmd/uid_provider/inprov/user_example.yaml`)
