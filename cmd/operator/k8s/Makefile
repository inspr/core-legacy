TAG :=$(shell git describe)

IMG_REPO:=gcr.io/red-inspr/inspr/

.ENVIRONMENT:=KAFKA_BOOTSTRAP_SERVERS \
INSPR_REGISTRY_URL \
INSPR_NODE_ID \
INSPR_NAMESPACE \
INSPR_LOG_CHANNEL \
INSPR_INPUT_CHANNELS \
INSPR_OUTPUT_CHANNELS \
KAFKA_AUTO_OFFSET_RESET \
KAFKA_ENABLE_AUTO_COMMIT

ENV := $(patsubst %,--env %="INSPR",$(.ENVIRONMENT))

MODULE_NAME:=k8s-operator

.PHONY: lint-docker
lint-docker:
	@docker build \
		--target lint \
		-t $(IMG_REPO)$(MODULE_NAME):$(TAG)-lint \
		.
	@docker run \
		--rm \
		$(IMG_REPO)$(MODULE_NAME):$(TAG)-lint \
		golint -set_exit_status ./... 

.PHONY: build-docker
build-docker:
	@docker build \
		--target build-env \
		-t $(IMG_REPO)$(MODULE_NAME):$(TAG) \
		.
.PHONY: push
push:
	@docker push $(IMG_REPO)$(MODULE_NAME):$(TAG)
		
.PHONY: test-docker
test-docker: 
	@docker run \
		--rm \
		$(ENV) \
		$(IMG_REPO)$(MODULE_NAME):$(TAG) \
		go test -cover ./...

.PHONY: build
build:
	@go build -o $(MODULE_NAME) main.go

.PHONY: test
test:
	@go test -cover ./...

.PHONY: lint
lint:
	@golint -set_exit_status ./...