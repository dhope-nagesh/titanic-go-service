# ==============================================================================
# Variables
# ==============================================================================

# Your container registry. Can be overridden, e.g., `make push REGISTRY=myregistry`
REGISTRY ?= nageshdhope

# Image names
SVC_IMAGE_NAME := titanic-go-service
DATA_IMAGE_NAME := titanic-data-container

# Image tag. Default to v1.0.0, can be overridden, e.g., `make build TAG=v1.0.1`
TAG ?= v1.0.0

# --- Environment Configuration ---
# Set the Kubernetes environment. Can be 'kind' or 'docker-desktop'.
K8S_ENV ?= docker-desktop

# Helm configuration
HELM_RELEASE_NAME := titanic-release
KIND_CLUSTER_NAME := titanic-cluster
DATA_SOURCE ?= csv # Default data source, can be overridden e.g., `make install DATA_SOURCE=sqlite`


# ==============================================================================
# Targets
# ==============================================================================

# Use .PHONY to ensure these targets run even if files with the same name exist.
.PHONY: all build push svc-image data-image setup install uninstall clean help test

# Default target runs when you just type `make`.
all: help

## --------------------------------------
## Setup & Dependencies
## --------------------------------------

install-swagger:
	@echo "--> Installing/updating Swagger CLI tool..."
	@go install github.com/swaggo/swag/cmd/swag@latest

generate-swagger: install-swagger
	@echo "--> Generating Swagger documentation..."
	@swag init -g cmd/server/main.go

seed-sqlite:
	@echo "--> Seeding SQLite database with data from CSV..."
	@go run cmd/seed/main.go


## --------------------------------------
## Image Building & Pushing
## --------------------------------------

build: svc-image data-image

svc-image: generate-swagger
	@echo "--> Building service image: $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG)"
	@docker build -f Dockerfile -t $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG) .

data-image: seed-sqlite
	@echo "--> Building data image: $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG)"
	@docker build -f Dockerfile.data -t $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG) .

push:
	@echo "--> Pushing images to registry: $(REGISTRY)"
	@docker push $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG)
	@docker push $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG)


## --------------------------------------
## Kubernetes Deployment (Kind & Helm)
## --------------------------------------

# This is the main setup target that chooses the correct environment.
setup: build
ifeq ($(K8S_ENV), kind)
	@if ! kind get clusters | grep -q "^$(KIND_CLUSTER_NAME)$$"; then \
		echo "--> Creating Kind cluster: $(KIND_CLUSTER_NAME)"; \
		kind create cluster --name $(KIND_CLUSTER_NAME); \
	else \
		echo "--> Kind cluster '$(KIND_CLUSTER_NAME)' already exists."; \
	fi
	@echo "--> Setting up kubectl context for Kind cluster: $(KIND_CLUSTER_NAME)"
	@kubectl config use-context kind-$(KIND_CLUSTER_NAME)
	@echo "--> Loading images into Kind cluster..."
	@kind load docker-image $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG) --name $(KIND_CLUSTER_NAME)
	@kind load docker-image $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG) --name $(KIND_CLUSTER_NAME)
else
	@echo "--> Using Docker Desktop Kubernetes. Ensure it is enabled and running."
	@kubectl config use-context docker-desktop
	@echo "--> Images built locally are automatically available. Skipping image load."
endif

# Main install target.
install: setup
	@echo "--> Installing/upgrading Helm release '$(HELM_RELEASE_NAME)' with data source '$(DATA_SOURCE)'"
	@helm upgrade --install $(HELM_RELEASE_NAME) ./helm/titanic-chart \
		--set image.repository=$(REGISTRY)/$(SVC_IMAGE_NAME) \
		--set image.tag=$(TAG) \
		--set dataLoaderImage.repository=$(REGISTRY)/$(DATA_IMAGE_NAME) \
		--set dataLoaderImage.tag=$(TAG) \
		--set config.dataSource=$(DATA_SOURCE)

uninstall:
	@echo "--> Uninstalling Helm release: $(HELM_RELEASE_NAME)"
	@helm uninstall $(HELM_RELEASE_NAME)


## --------------------------------------
## Testing & Cleanup
## --------------------------------------

test:
	@echo "--> Running tests..."
	@go test ./... -v

clean:
ifeq ($(K8S_ENV), kind)
	@echo "--> Deleting Kind cluster: $(KIND_CLUSTER_NAME)"
	@kind delete cluster --name $(KIND_CLUSTER_NAME) --quiet
else
	@echo "--> To clean Docker Desktop, reset the Kubernetes cluster in the Docker Desktop UI."
endif

help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build both the service and data docker images."
	@echo "  push          Push docker images to the configured registry."
	@echo "  install       Build images, set up K8s (if kind), and deploy the application."
	@echo "  uninstall     Remove the application from the K8s cluster."
	@echo "  test          Run all Go tests."
	@echo "  clean         Delete the local Kind cluster (if using kind)."
	@echo "  help          Show this help message."
	@echo ""
	@echo "You can override variables, e.g.:"
	@echo "  make install K8S_ENV=docker-desktop"
	@echo "  make install TAG=v1.0.2 REGISTRY=myrepo DATA_SOURCE=sqlite"
	@echo ""
