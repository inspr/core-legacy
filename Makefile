# vim: set foldmarker={,} foldlevel=0 foldmethod=marker:

-include ./.env

VERSION ?= $(shell git describe --long --always --dirty)
PROFILE ?= inspr-stack
K8S_NAMESPACE ?= ${VERSION}
GOFLAGS ?= 
RELEASE_NAME ?= ""

VALUES ?= stack-overwrites.yaml
INSPRD_VALUES ?= insprd-overwrites.yaml
UIDP_VALUES ?= uidp-overwrites.yaml
export



help:
	@python3 ./scripts/makefile_help.py Makefile $C

# golang {

## downloads dependencies
go/download:
	go mod download

## builds all binaries to the bin directory
go/build:
	mkdir -p build/bin
	go build ${GOFLAGS} -o build/bin ./...

## runs all tests on the repo
go/test:
	go test ./...

## runs all tests tagged integration
go/test/integration:
	go test ./... -tags=integration

## lints the repo with goling and staticcheck
go/lint:
	staticcheck ./...
	golint ./...

## gets tools necessary for development
go/init:
	go get -u golang.org/x/lint/golint
	go get -u github.com/ptcar2009/ptcwatcher/cmd/ptcwatcher


## runs coverage and exports coverage profile
go/coverage:
	bash ./.github/scripts/unittest.sh

## watches the files and lints on changes
go/lint/watch:
	ptcwatcher 'make go/lint' -w ./pkg -w ./cmd

## watches the files and builds on changes
go/build/watch:
	ptcwatcher 'make go/build' -w ./pkg -w ./cmd

## watches the files and tests on changes
go/test/watch:
	ptcwatcher 'make go/test' -w ./pkg -w ./cmd
# }

# CLI {

# insprctl {
## builds insprctl to the bin directory
cli/insprctl/build:
	go build -o bin ./cmd/insprctl

## installs insprctl to $GOPATH/bin
cli/insprctl/install:
	go install ./cmd/insprctl
# }

# inprov {
## builds inprov to the bin directory
cli/inprov/build:
	go build -o bin ./cmd/uid_provider/inprov

## installs inprov to $GOPATH/bin
cli/inprov/install:
	go install ./cmd/uid_provider/inprov
# }

# all {
## builds all CLI tools to the bin directory
cli/build: cli/insprctl/build cli/inprov/build

## installs all CLI tools to $GOPATH/bin
cli/install: cli/insprctl/install cli/inprov/install
# }
# }

# CI {

## runs all scripts regarding CI, including linting, coverage and initialization
ci/all:  ci/init ci/lint ci/test ci/coverage

## initializes the environment for CI
ci/init: go/init helm/init semgrep/init

## lints the go src, helm templates and runs semgrep jobs
ci/lint: go/lint semgrep/run helm/lint

## runs all tests regarding golang
ci/test: go/test go/test/integration

## runs coverage on the repo and exports the profile
ci/coverage: go/coverage

## builds the CLI to all platforms and syncs to the repo
ci/release: ci/cli/push

# CLI {
## builds the CLI to all platforms and architectures
ci/cli/build:
	bash ./.github/scripts/buildcli.sh

## pushes the built binaries to the CI repo
ci/cli/push: ci/cli/build
	bash ./.github/scripts/pushcli.sh
# }
# }

# helm {

# uidp {
## packages the UIDP helm chart using the UIDP overrides file, which by default is uidp-overwrites.yaml
helm/uidp/package:
	helm package ./build/uidp

## lints the UIDP helm chart using the UIDP overrides file.
helm/uidp/lint:
	helm lint ./build/uidp

## runs the UIDP helm chart tests using the UIDP overrides file,
helm/uidp/test:
	helm test ${RELEASE_NAME}-uidp -n ${K8S_NAMESPACE}

## installs the uidp helm chart to the K8S_NAMESPACE using the uidp overrides file.
helm/uidp/install:
	helm upgrade -i ${RELEASE_NAME}-uidp ./build/uidp -f ${UIDP_VALUES} -n ${K8S_NAMESPACE}

## delete a helm release
helm/uidp/uninstall:
	helm uninstall ${RELEASE_NAME}-uidp -n ${K8S_NAMESPACE}
# }

# insprd {
## packages the INSPRD helm chart using the INSPRD overrides file, which by default is insprd-overwrites.yaml
helm/insprd/package:
	helm package ./build/insprd

## lints the INSPRD helm chart using the INSPRD overrides file.
helm/insprd/lint:
	helm lint ./build/insprd

## runs the INSPRD helm chart tests using the INSPRD overrides file,
helm/insprd/test:
	helm test ${RELEASE_NAME}-insprd -n ${K8S_NAMESPACE}

## installs the insprd helm chart to the K8S_NAMESPACE using the insprd overrides file.
helm/insprd/install:
	helm upgrade -i ${RELEASE_NAME}-insprd ./build/insprd -f ${INSPRD_VALUES} -n ${K8S_NAMESPACE}

## delete a helm release
helm/insprd/uninstall:
	helm uninstall ${RELEASE_NAME}-insprd -n ${K8S_NAMESPACE}
# }

# stack {
## packages the INSPR-STACK helm chart using the INSPR-STACK overrides file, which by default is stack-overwrites.yaml
helm/package:
	helm package ./build/inspr-stack

## lints the INSPR-STACK helm chart using the INSPR-STACK overrides file.
helm/lint:
	helm lint ./build/inspr-stack

## runs the INSPR-STACK helm chart tests using the INSPR-STACK overrides file,
helm/test:
	helm test ${RELEASE_NAME}-stack -n ${K8S_NAMESPACE}

## installs the inspr-stack helm chart to the K8S_NAMESPACE using the inspr-stack overrides file.
helm/install:
	helm upgrade -i ${RELEASE_NAME}-stack ./build/inspr-stack -f ${VALUES} -n ${K8S_NAMESPACE}

## delete a helm release
helm/uninstall:
	helm uninstall ${RELEASE_NAME}-stack -n ${K8S_NAMESPACE}
# }
# }

# Skaffold {

## runs skaffold build with the PROFILE profile and outputs the image to OUTPUT_FILE if defined.
skaffold/build:
ifdef OUTPUT_FILE
	skaffold build -p ${PROFILE} -o ${OUTPUT_FILE}
else
	skaffold build -p ${PROFILE}
endif

## runs skaffold run with the PROFILE profile on the K8S_NAMESPACE namespace.
skaffold/run:
	skaffold run -p ${PROFILE} -n ${K8S_NAMESPACE}

skaffold/dev:
	skaffold dev -p ${PROFILE} -n ${K8S_NAMESPACE}

## Deletes the release and the namespace that it was created in
skaffold/delete:
	skaffold delete -p ${PROFILE} -n ${K8S_NAMESPACE}
	kubectl delete namespace ${K8S_NAMESPACE}
# }

# semgrep {
## downloads sempgrep and installs it using python3
semgrep/init:
	python3 -m pip install semgrep

## runs the desired test suites for semgrep
semgrep/run:
	semgrep --config "p/trailofbits"
# }

# secrets {
## gets and decodes the admin secret
secrets/uidp/admin:
	@echo $(shell kubectl get secrets -n ${K8S_NAMESPACE} ${RELEASE_NAME}-init-secret -o jsonpath="{.data.ADMIN_PASSWORD}" | base64 --decode)

## gets and decodes the insprd init key secret
secrets/insprd/init:
	@echo $(shell kubectl get secrets -n ${K8S_NAMESPACE} ${RELEASE_NAME}-insprd-init-key -o jsonpath="{.data.key}" | base64 --decode)

## gets and decodes the grafana admin password
secrets/grafana/password:
	@echo $(shell kubectl get secrets -n ${K8S_NAMESPACE} ${RELEASE_NAME}-grafana-admin -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" | base64 --decode)
# }

# dashboards {
## port forwards grafana and opens a browser session with it
dashboards/grafana:
	xdg-open http://localhost:3000
	kubectl port-forward -n ${K8S_NAMESPACE} $(shell kubectl get pods --namespace ${K8S_NAMESPACE} -l "app.kubernetes.io/name=grafana" -o jsonpath="{.items[0].metadata.name}") 3000:3000

## port forwards prometheus and opens a browser session on it
dashboards/prometheus:
	xdg-open http://localhost:9090
	kubectl port-forward -n ${K8S_NAMESPACE} $(shell kubectl get pods --namespace ${K8S_NAMESPACE} -l "app.kubernetes.io/name=${RELEASE_NAME}-prometheus" -o jsonpath="{.items[0].metadata.name}") 9090:9090
# }

# port forwards {
pf/insprd:
	kubectl port-forward -n ${K8S_NAMESPACE} $(shell kubectl get pods --namespace ${K8S_NAMESPACE} -l "app.kubernetes.io/name=insprd" -o jsonpath="{.items[0].metadata.name}") 8080:8080
pf/auth:
	kubectl port-forward -n ${K8S_NAMESPACE} $(shell kubectl get pods --namespace ${K8S_NAMESPACE} -l "app.kubernetes.io/name=auth" -o jsonpath="{.items[0].metadata.name}") 8081:8081
pf/uidp:
	kubectl port-forward -n ${K8S_NAMESPACE} $(shell kubectl get pods --namespace ${K8S_NAMESPACE} -l "app.kubernetes.io/name=uidp" -o jsonpath="{.items[0].metadata.name}") 9001:9001
# }


