.PHONY: all
all: build

.EXPORT_ALL_VARIABLES:
PATH = $(LOCALBIN):$(shell echo $$PATH)

.PHONY: build
build: fmt vet ## Build manager binary.
	@mkdir -p $(LOCALBIN)
	go build -o $(LOCALBIN)/batchv1-controller main.go

run:
	./.build/batchv1-controller --debug

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

LOCALBIN ?= $(shell pwd)/.build
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CONTROLLER_GEN := $(GOPATH)/bin/controller-gen
$(CONTROLLER_GEN):
	pushd /tmp; go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.3; popd

crd-manifests: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) crd:maxDescLen=0 paths="./pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1/..." output:crd:artifacts:config=crds

.PHONY: generate
generate: $(CONTROLLER_GEN) codegen-clientset crd-manifests

codegen-clientset:
	@echo "Generating Kubernetes Clients"
	./hack/update-codegen.sh

ensure-go-junit-report:
	@command -v go-junit-report || (cd /tmp && go install github.com/jstemmer/go-junit-report/v2@latest)

test: ensure-go-junit-report
	export PATH=$$PATH:~/go/bin:$$GOROOT/bin:$$(pwd)/.build; go test -v ./... -covermode=count -coverprofile=coverage.out 2>&1 | go-junit-report -set-exit-code -out junit.xml -iocopy
