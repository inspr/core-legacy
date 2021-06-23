
# Changelog


### #79 Story CORE-481 | Adding fields descriptions to the documentations in docs/yamls 
- misc:
    - created the alias doc
    - added fields to the node documentation
    - checked index and added the alias section to it
---

### #78 Bug CORE-475 | Deletion of multiple channels not working  
- fixes:
    - Altered the Channel Operator when handling delete to get the channel from `tree` instead of `root`
    - Fixed Insprd CLI's Alias delete command (and its tests)
- refactors:
    - In the tree memory manager renamed method Root(), which returns the tree, to Tree()
    - In the tree memory manager, renamed "RootGetter" structures to "PermTreeGetter"
- misc:
    - Updated Insprd CLI's cluster config command error messages
---

## Release v0.1.0

### #76 Tech CORE-360 | vanity url
- fixes:
  - changed the `module` name in the `go mod` of the project, that allows the vanity url in the productions cluster to obtain the package properly
---

### #77 Story CORE-463 | Create documentation for Auth
- features:
    - authentication doc
---

### #71 Story CORE-462 | Update old docs
- features:
    - Updated most of the files in /docs folder
    - Updated scripts that build and install Inspr CLI
    - Updated release cretion GitHub Action so it pushes the CLI once it's build
- fixes:
    - Fixed Ping Pong example
    - Fixed resolved channel env var injection in converter.go so it injects the original channel to reference the resolved one
- misc:
    - Renamed folder from apply_channel_type.go to apply_type.go in `cmd/insprctl/`
    - Added/updated some logs, error messages and method comments all around
---

### #65 Story CORE-401 | CLI create new broker command
- features:
    - developed the inspr cli subcommand for reading a yaml file and trying to create a broker in the cluster
    - created the interface method in the controller interface
    - implemented the create broker method
    - added the `/brokers/kafka` endpoint in the insprd routes
    - added the handler for the `/brokers/kafka` endpoint
- fixes:
    - added the simple test of the `/brokers` routes.go in the api/controller pkg
    - changed from `inspr` to `insprctl` in docs that referenced the CLI command
- refactors:
    - removed the `parser` package in the directory `meta/utils/parser`, all of its contents can now be found on the `yaml_parser.go` file in the `meta/utils` pkg.
    - removed the []byte to []byte comparison in the inspr cli testing since it can be just a direct comparison between two strings. No need for adding extra buffers 
---
### #75 Story CORE-467 | Modify operator structure to handle multibrokers
- features:
   - configured broker memory structure do store configurations
   - configured operator manager to import configurations dynamically so that operator always work with broker configured
- fix:
    - fixed tests
    - removed unnecessary connection port from broker-specific sidecars
---

### #74 Tech CORE-430 | Refactor memory dApp section 
- refactors:
    - Update error messages so they use Sprintf structure and aren't duplicated in `dapp_utils.go`
    - Methods in `dapp_utils.go` now use multi-error structure
    - Removed unused methods in `dapp_utils.go`
- fixes
    - Added LB Sidecar, Authsvc and Secretgen images tag update in Helm Chart update script
---

### #73 Tech CORE-461 | Rename inspr CLI as insprctl
- misc:
    - Renamed folder `cmd/inspr` to `cmd/insprctl`
    - Modified commands and command descriptions in `pkg/cmd` and `cmd/insprctl`
    - Fixed tests that were broken after the main command change
---

### #72 Tech CORE-400 | Revaluate nomenclature of scope
- refactors:
    - Renamed comments, func params and variables from "context" to "scope"
    - Refactored some error messages (better messages and use sprintf)
    - Renamed env var "INSPR_APP_CTX" to "INSPR_APP_SCOPE"
    - Moved "InitDO" to models file in pkg/auth
---

### #70 Tech CORE-442 | Inspr CLI describe command outdated
- features:
    - PrintAppTree function now prints the Node Meta, the Node Ports, the Node Sidecar LBRead and LBWrite port and the new Auth structure.
    - PrintChannelTree function now prints the ConnectedAliases, the PriorityBrokerList and the SelectedBroker
---

### #69 Story CORE-416 | k8s operator different sidecar injection
- features:
    - operator now deploys a node, its LB sidecar and all of its broker-specific sidecars
    - updated Insprd's broker creation method
    - updated how kafka sidecar is configured to k8s deployment
    - updated Sidecar Factory so now it can receive more container options, and it returns that sidecar's environment variables (such as its port and address)
    - created utils method that returns free/available tcp ports
    - added LB sidecar to Helm Chart and Skaffold
- fixes:
    - removed old sidecar structures
    - removed kafka configurations from Helm Chart and Skaffold
    - removed odd structure that was being return by sidecar client channel handler
    - fixed some comments, logs and funcion names misspells
    - fixed LB Sidecars handlers that weren't getting the schema from the resolved channel
    - fixed broker-specific sidecar read handler (it wasn't using a channel as path when sending request to lb sidecar)
- refactors:
    - moved brokers constants from `cmd/insprd/memory/brokers` into `pkg/meta/brokers`
- misc:
    - created new example to test multibroker architecture using Kafka
---

### #66 Story CORE-415 | Kafka Sidecar redesign
- features:
    - sidecar generic structure reconfigured to support broker specific sidecars
    - new kafka sidecar created
    - kafka sidecar image created for broker specific structure
    - node operator reconfigured to work with broker specific sidecar
    - environment methods reconfigured to fetch node data in new format
- fixes:
    - rest error unmarshaler method fixed and tested
---

### #68 Bugfix CORE-431 | UIDP Refresh token crash
- refactors:
    - Improved some error's messages, most of them in `cmd/authsvc`
    - Changed some logs so they use `zap.logger` instead of `log` and improved their messages
    - Fixed misspelings pointed out by Go Report
    - Changed "POST", "GET", "DELETE" and "PUT" to use `net/http` constants instead
    - Renamed the embed file used in `examples/controller`
---

### #67 Tech CORE-429 | Change the inspr cmd cluster command to ignore the .inspr/token
- features:
    - Updated the Init function in the auth package so now it creates .inspr/token if it doesn't exist
    - Added the admin permissions in auth models
    - The admin user is now initialized with the admin permissions created in auth models
---

### #63 Story CORE-414 | Adapt the current sidecar to a load balancer sidecar
- features:
    - implemented the new LB sidecar
    - created new functions to retrieve sidecar's env vars
    - updated the converter to generate new env vars used by the LB sidecar
    - changed the env vars associated with the old sidecar so they are now used by the LB sidecar
- refactors:
    - renamed the old sidecar structure to `sidecar_old` and fixed where it was imported. This structure will be maintained while the multibroker sidecar isn't released
    - renamed some variables/parameters so they make more sense
---

### #60 Story CORE-419 | Create Kafka sidecar Factory
- features:
    - developed a function that returns a `SidecarFactory` type, responsible for
      the sidecar deployment configuration in the k8s.
    - opted for the implementation of the sidecar factory to be in the
      `sidecars` pkg.
    - created the utils file containing small functions and mocks of functions
      to be, these functions should be able to be used by other sidecars.
---

### # 64 Tech CORE-448 | Unification of insprctl cluster commands
- fixes:
    - tested cluster commands `brokers` and `init`
- refactors:
    - changed `cluster` command to follow pattern of inspr commands
    - changed WithAliases method from command builder to receive a variadic parameter
---

### #62 Story CORE-413 | Modify Channel structure
- features
    - Added SelectedBroker and BrokerPriorityList to Channel structure
    - Create a function that sets the SelectedBroker of a Channel based on the BrokerPriorityList and the available brokers
    - Call the function above in create App and Channel
    - Added a function in brokers that resets the manager (for tests)
---

### #61 Story CORE-428 | CLI get installed brokers command
- features:
    - added a Brokers method and Brokers interface to controllers interface
    - changed controller client to comply with new interface
    - implemented a broker data interface (DI)
    - changed insprd's server struct, http client, handler structure, and their mocks to include a BrokerManager
    - created a handler to get broker data from insprd
    - implemented cli command for getting broker data
-  tests:
    - unit tested handler
    - integration test succeeded on cluster
---

### #58  Story CORE-418 | Create Sidecar Factory
- features:
    - developed a function type to represent sidecar factories
    - declared a interface for sidecars and implemented its structure as an abstract factory
    - developed a Manager interface to envelop BrokerManager and  AbstractBrokerFactory
    - implemented Manager's struct
    - implemented SidecarInterface's methods
    - tested previously mentioned methods
- fixes:
    - moved broker management interfaces and structures to a more appropriate directory
---

### #57 Tech CORE-426 | Review Permissions
- features:
    - created functions in the UIDP client that make sure when creating a token, it’s permissions should be the same as the creator’s permissions (or have less permissions).
- fixes:
    - all functions of the get command. For example, if you're trying to get channels, before calling the getDapp directly it will call the getChannels to check if the user have the right permissions.
---

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
