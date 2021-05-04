
# Changelog

### #43 Issue CORE-335 | Controllers configuration for in cluster initialization.
- fixes:
        - Renamed poorly named rest clients on controller clients.
- features:
	- ControllerConfig object developed to help initialize controllers and ease of config validation.
	- GetInClusterConfigs function implemented to fetch controller configs from current in-cluster dapp deployment
- tests:
	- unit tests for GetInClusterConfigs
---

### #40 Issue CORE-332 | adds a permissions field to the dApp definition

- features:
	- adds a permission field to dapp meta
	- changes new dapp on tree to handle the injection
---


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
