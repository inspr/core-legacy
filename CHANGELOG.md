# Changelog

### #133 Feature Grafana dashboards | Automatically create becnhmark dashboard on grafana deploy
- feature:
  - Added support for custom dashboards in grafana values
  - Added dashboard config map
  - Moved dashboard file to the helm chart
---

### #130 Story CORE-575 | Develop metrics for the sidecar broker Reader
- feature:
  - added metrics to the kafka reader structure
  - created GetMetric function inside reader writer
  - added metric for the resolve channel duration
---

### #131 Hot-Fix (B.M example): benchmarking dapp
- fix:
  - changed default bench dapp scope 
---

### #129 Story CORE-574 | Develop metrics for the sidecar broker Writer
- feature:
  - added metrics to the kafka writer structure
  - created GetMetrics function inside kafka writer
  - added metric for the resolve channel duration
  - added metric for the produce message duration
  - added metric for the flush duration
---

### #127 Story CORE-572 | Permission scheme inversion
- features:
  - refactored auth midleware
  - inverted permission scheme
  - adapted validation methods
  - refactored cluster init flow
  - removed unused scope headers
--- 

### #128 Story CORE-563 | Develop Sidecar metrics
- feature:
  - added two Prometheus Summary metrics for the sidecar
---

### #125 Story CORE-567 | Create custom queries for Inspr's Prometheus metrics
- doc:
  - added the json for the internal confluence documentation
- misc:
  - added a deployScript in the `bench_mark` example
---

### #117 Story | Inspr stack workspace initialization
- doc:
	- create helm installation guides to insprd, uidp and inspr-stack, with some overwrite examples
	- update the overwrite values configuration table
- dev:
  - update helm chart Makefile commands and paths
---

### #124 Story CORE-565 | Develop sidecar throughput metrics
- feature:
  - created constants for the channels names in the ping pong example
  - added a script to deploy the ping pong
---

### #123 prometheus chart selector labels setting incorrect name
- fix:
  - Makefile using wrong variable in secrets commands
  - prometheus labels were being set to commonLabels, which was setting k8s.app to
    releaseName,this was changed to ${releaseName}-prometheus
  - Makefile dashboard/prometheus now works with the helmChart default values, doens't work if
    someone sets `values.prometheus.fullname` as this overwrites the template generated values.
---

### #122 Story CORE-566 | Implement benchmarking dapp
- feature:
    - throughput bench marking dapp
- fix:
    - make target for getting grafana password
---

### #121 Quick gitgnore fix

- fix:
  - set up to not track new helm overwrites files

---

### #109 Prometheus setup

- features:
  - configured prometheus and grafana charts on inspr-stack chart
- fix:
  - renamed uidp directory

---

### #120 Story CORE-502 | Develop client (for controller dApps) for Python

- features:
  - models module added to inspr python lib (to define the changelog and inspr structure data types)
  - rest module updated (added get, delete and update methods)
  - controller-client module added to inspr python lib (with all the CRUD methods for apps, channels, types and alias)

---

### #118 Story CORE-561 | Authentication foir broker commands

- features:
  - added header scope on broker requests
  - configured validation for broker operations

---

### #113 Story CORE-560 | CLI changes documentation

- doc:
  - Updated docs and examples to host flag and Brokers command

---

### #111 Feature | adds redis as a dependency for the uidp helm chart

- features:
  - adds redis as a uidp helm dependency

---

### #112 Bug CORE-569 | Cli initialization not adding "http://" to ip

- fix:
  - added http:// on the `serverip` on the `insprctl init` command

---

### #110 Story CORE-246 | Develop a client (for sidecar communication) with Python

- features:
  - Created the newClient, writeMessage, handleChannel and run functions for the lb sidecar client for python
  - Created a ping pong example

---

### #108 Tech CORE-549 | Add support to request Host header in the cli

- features:
  - added host to the client structure in the request and controller pkg, to direct serverip requests
  - created host flag `--host` to overwrite host configuration on insprctl commands without changing the host configured on the client
- tests:
  - updated the set client argument to fit on the new host request call
  - add on request_test.go `host` arg and `test_new_host` case

---

### #105 Tech CORE-442 | Refactor of the error structure in the project

- features:
  - refactor of the ierror pkg, removal of the builder concept, now has
    similar behaviour to the standard library errors pkg but functions to add
    more context to the error.
  - The main focus was to develop a secluded pkg to be used for testing and
    error context, should be used as a referenced in future features.
- tests:
  - changed 80+ files tests
  - redid the tests of the ierror pkg

---

### #107 Bug CORE-550 | Mutex not unlocking on panics

- fix:
  - panic recovery methods now receive a cancellation function that, when not null, is executed when panics are recovered.

---

### #106 Tech CORE-492 | Correction on the controller operation

- misc:
  - improved debugging messages for dapp creation and deletion
- fix:
  - created a controll variable on k8s operator that alternates between permanent and changed memory throughout the the entire operation on nodes
  - Added log level controll to authservice

---

### #100 Feat/Fix Pprof/memLeak | adding pprof routes to servers and fixing memory leak

- fix:
  - closed the http.Requests and http.Response
- features:
  - added the pprof routes in the loadbalancer and services

---

### #101 Fix | adds routine for collecting kafka events

- fix:
  - update kafka go version to 1.7.0
  - adds go routine for collecting kafka writer events

---

### #103 Tech CORE-542 | Change broker commands and sub-commands

- features:
  - change clusters cli sub-commands.
  - create brokers cli command and sub-commands.
  - update workspace_init doc to the new broker command.
- tests:
  - fixed the tests that are related to brokers and clusters.

---

### #104 Tech CORE-545 | external-dns and removal of hostname in cli

- fix:
  - removed the hostname in the insprctl requests
  - added small documentation about external-dns
- misc:
  - added `external-dns-*.yaml` to gitignore

---

### #102 Devops | Updated script for chart releassing

- features:
  - updated github integration scripts for helm chart release updating

---

### #101 Fix | adds routine for collecting kafka events

- fix:
  - update kafka go version to 1.7.0
  - adds go routine for collecting kafka writer events

### #98 Story CORE-493 | Restart Policy for Nodes

- refactors:
  - Deleted option/function of configuring restartPolicy of a deployment
  - Deleted the tests of the function that changed the restart policy
  - Added support for creating nodes with replicas

---

### #97 Tech CORE-490 | Channels envvars with broker

- features:
  - change the separator for the channel brokers.
- test
  - fixed the tests that use channel's env variables

---

### #95 Tech CORE-491 | Improve UIDP

- features:
  - added a logger to the requests that the uidp handles.
  - added description for the cli commands in the `inprov`.
- fixes:
  - removed the `COMMIT SUCCEEDED` message when the git hook didn't encounter any errors.

---

### #96 Tech CORE-372 | Protect user password on inprov login #96

- features:
  - developed new execution standard for login command
- tests:
  - adapted tests

---

### #93 Tech CORE-366 | Encrypt the user password into the database

- features:
  - Password encryption when creating user
  - Password encryption when initializing admin user
  - Passord verifycation when trying to delete/update/create a new user now uses the hashed password and bcrypt function
  - Redis install/uninstall documentation
- refactors:
  - Updated all the tests for create/delete/update a user

---

### #91 Bug CORE-498 | Unable to update channels

- fixes:
  - Changed the Annotations metadata diff to work according to new brokers configuration process.
  - Copied `SelectedBroker` from old channel to new on Update operation
- misc:
  - changed Skaffod deployment to get secret values from a file other than the default Values.

---

### #92 Tech CORE-456 | Review auto labeler

- features:
  - updated autolabeler workflow so it labels properly
  - updated codecov action to use codecov's github action

---

### #90 Bug CORE-499 | Unable to delete multi-layered dApps

- refactors:
  - Added the usePermTree (bool) parameter in the ResolveBoundary function. If true, then use the unmodified app tree.
  - Updated all uses of the ResolveBoundary function. In particular, in the convert file of package node the parameter usePermTree is set to true

---

### #89 Story CORE-500 | removing inspr.com from local `/etc/hosts`

- fixes:
  - added hostname variable to the Send func in the request pkg
  - changed the git hook since `staticcheck` and `go vet` fails when looking only at a specific file, now looks at the entire repository folder
  - changing minikube.md to replace `/etc/hosts` with a explanation on how to setup the serverip via the insprctl cli.

---

## Release v0.1.1

### #87 Story CORE-478 | readme.md

- features:
  - updated readme.md doc

---

### #86 Story CORE-485 | minikube.md

- features:
  - created the documentation on how to run insprd/uidp on minikube

---

### #85 Story CORE-464 | Multibroker documentation

- features:
  - documentation for multibroker feature

---

### #84 Story CORE-480 | difference.md

- features:
  - Updated `difference.md`
- refactors:
  - Changed parameter from `ctx` to `scope` in `pkg/meta/utils/diff/diff.go`

---

### #81 Story CORE-476 | Memory dependency refactor \[Singletons\]

- features:
  - design new general memory manager to provide a single access point to all memory for Insprd
  - implemented new general MemoryManager by injecting both tree and brokers managers into it.
  - removed all broker dependency from tree managing structure
  - injected necessary broker information into tree managing metheds that require it
  - refactored singletons to allocate managing structs instead of interfaces to better fit the pattern
  - removed all access points to tree memory manager through the singleton from the tree managing methods, replaced it's uses to the access point that already existed inside each structure specific manager
  - improved structure specific permanent tree getters to be capable of connecting back to the general permanent tree getter
  - removed all access points to tree memory manager through the singleton from the structure specific permanent tree getters
- fixes:
  - fixed all tests
  - fixed k8s operator example to comply with new data injected methods
- misc:
  - added multiplexers to protect broker data structure

---

### #83 Story CORE-477 | troubleshooting.md

- features:
  - Updated `troubleshooting.md` with some new "troubles" that users may encounter
  - Updated PingPong demo readme file

---

### #82 Story CORE-482 | fix examples

- features:
  - Documented all examples/demos in `/examples`
  - Updated all structures YAML definitions for the examples/demos
  - Updated all Dockerfiles and `.go` files for examples/demos
  - Updated PubSub doc

---

### #50 - Feature: creates git hooks scripts and pre-commit

- feature:
  - creates a wrapper script for git hooks
  - creates a pre-commit script for linting files that are staged for committing
  - creates an installation script that creates symlinks for the hooks
  - added a readme describing the process of installation
    - added a comment that signals that the staticcheck should ignore the file in:
      - `node_operator.go`
      - `node_operator_test.go`

---

### #80 Story CORE-497 | Migrate cluster from "red-inspr" to "insprlabs"

- misc:
  - changed the tags of the dockers image, from red-inspr to insprlabs. Meaning that i changed the place in which the docker images are stored.

---

### #79 Story CORE-481 | Adding fields descriptions to the documentations in docs/yamls

- misc:
  - created the alias doc
  - added fields to the node documentation
  - checked index and added the alias section to it

---

### #33 update networking from beta1 to v1 in the ingresses

- fixes:
  - Altered the fields of all the ingresses that used the previous version of the v1beta1 of k8s.
  - `authsvc/secret` now passes context in the methods

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
- tests:
  - unit tested handler
  - integration test succeeded on cluster

---

### #58 Story CORE-418 | Create Sidecar Factory

- features:
  - developed a function type to represent sidecar factories
  - declared a interface for sidecars and implemented its structure as an abstract factory
  - developed a Manager interface to envelop BrokerManager and AbstractBrokerFactory
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

### #40 Issue CORE-332 | adds a permissions field to the dApp definition

- features:
  - adds a permission field to dapp meta
  - changes new dapp on tree to handle the injection

---

### #50 - Feature: creates the completion command on the CLI

- feature:
  - creates the completion command
- fix:
  - adds descriptions to the cluster init command

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
  - changed Payload (`pkg/auth/models/payload.go`) field 'Role' (`int`) to 'Permission' (`map[string][]string`), and removed field 'Scope' - 'Permission' is now a map of Scope to list of Permissions
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
- tests: - changes kafka tests to handle the new channel based approach
  =======
