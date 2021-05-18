
# Changelog

### #59 - Fix: CORE-310 | Automatic Helm chart tags update
- fixes:
    - The workflow/GitHub action `release.yml` wasn't properly working. To fix it:
        - Added action to setup Skaffold and Helm tools
        - Added action to setup Cloud SDK and configure the service account
---

### #56 - Story CORE-417 | Create a broker-tracking structure
- features:
    - created a Brokers struct on meta/brokers to store broker information.
    - created a Brokers interface to handle storage of broker data
    - included said interface on MemoryManager
    - developed a mock structure for the new interface
---

### #55 - Tech CORE-310 | Automatic Helm chart tags update
- features:
    - created new GitHub action to automatically create a new release if a tag is pushed
    - Inspr CLI binaries are now going to be available in the GitHub Release instead of GCloud
---

### #54 - Tech CORE-365 | Alias Meta.Name automatic injection when applying app
- features:
    - alias name injection on apply app
    - updated user example to include token creation
    - created a alias naming structure validator
    - created a alias parser from compound scope
    - created a cli argument processor for alias names
- fixes:
    - describe alias cli command
- refactors:
    - removed code duplication from cli's dApp apply functionality
---

### #53 - Tech CORE-362 | Validate alias when applying app
  - features:
    - created a validation function for aliases and used it on Create and Delete
---
    
### #52 - Bug CORE-425 | Check failing controller update method
- fixes:
    - fixed the Controller dApp example so create/update/delete methods are working
- misc:
    - formated the CHANGELOG.md file so all changelogs have the same format
---

### #48 - Feature: adds init command to the cli
- features:
    - creates an init command for configuration initialization
    - adds config flag to specify configuration file
    - changes flow so that if the config is invalid, asks for initialization
- fixes:
    - silences main command to remove duplicate errors
- refactors:
    - refactors configuration initialization to allow different files
    - changes command builder to now use variadic flag argument instead of slice
---

### 46 - misc: changing the name of channelType to type
- fixes:
    - use the replace funcionality in vscode to replace `Channel_Type` for `Type`
    - use the replace funcionality in vscode to replace `ChannelType` for `Type`
    - use the replace funcionality in vscode to replace `Channel Type` for `Type`
    - changed the DI structures used in the cli handlers
    - changed the tests of the handlers, updated the query being used to test
      the methods
    - changed a few file names for it to have just `type` instead of
      `channel_type`
    - left some places with `type` as is more intuitive as to what type is
      being managed.

### #44 - Issue CORE-333 | Create the token for authentication inside the cluster
- features:
	- creates a secret for each node with its permissions
	- adds environment variables necessary for controller authentication
- refactors:
	- major refactor of the k8s operator conversion system
	- creates a helper package for kubernetes deployment and container building
	- changes node operator operations to be interface based
---

### #43 - Issue CORE-335 | Controllers configuration for in cluster initialization.
- features:
	- ControllerConfig object developed to help initialize controllers and ease of config validation.
	- GetInClusterConfigs function implemented to fetch controller configs from current in-cluster dapp deployment
- fixes:
    - Renamed poorly named rest clients on controller clients.
- tests:
	- unit tests for GetInClusterConfigs
---

### #41 - Feature: Handling invalid requests in the middleware authorization
- features:
    - rest/middleware.go -> using the `permissions` field in the payload, created in the issue CORE-336, for the scope and permissions check.
    - reading the scope information from the http.request header, instead of the body
    - changing the controller/client so is more functional in nature, functions return a copy of the struct with the new value.
    - rest/request -> changed so the scope is inserted in the header of the request, the structs that were the body of the request were changed so they no longe have a `scope` field.
- fixes:
    - fixed tests in the controller/client and cli/handlers that were broken due to some struct changes.
- tests:
    - new tests for the middleware, since there is a new logic step due to the permission introduction.
---

### #40 - Issue CORE-332 | adds a permissions field to the dApp definition
- features:
	- adds a permission field to dapp meta
	- changes new dapp on tree to handle the injection
---

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
- refactors:
    - grouped files under `pkg/auth/models` into `models.go` in `pkg/auth`, and deleted the folder (this was done so the structures have the package prefix `auth` when being imported)
        - Adapted all methods which use the structures defined in `models` so they continue to work as intended
- tests:
	- changes in tests of methods which were impacted by Payload and User modification
- misc:
    - removed `init()` func that wasn't working in `cmd/insprd/operators/kafka/nodes/converter.go` (@ptcar2009 fixed it PR #37)
    - created example yaml file for uid provider user's structure definition (`cmd/uid_provider/inprov/user_example.yaml`)
---

### #37 - Feature: reactive sidecar <!-- This is the title -->
- features:
	- sidecar server now has a goroutine for constant message polling. This routine reads from all subscribed channels and sends requests to the main dApp so that the dapp handles those requests
	- sidecar's read message now has a maximum number of retries
	- sidecar and dApp now communicate via tcp instead of unix sockets. This is helpful to add easier integration with libs for many languages, including requests for js, Flask for python, etc. 
	- inspr client now acts as a wrapper on a http server for reading messages
	- write message now takes the routes as channel parameters
	- adds fields for port definition on the dApp metadata
- fixes:
	- lowers the timeout on kafka writer's flush
	- changes some loggers to be production and lower test cluttering
	- some type fixes
- tests:
	- changes kafka tests to handle the new channel based approach
=======
