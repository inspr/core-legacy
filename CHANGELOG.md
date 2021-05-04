# Changelog

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