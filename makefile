# ==============================================================================
# Variables
# ==============================================================================

# Your container registry. Can be overridden, e.g., `make push REGISTRY=myregistry`
REGISTRY ?= yourdockerhubusername

# Image names
SVC_IMAGE_NAME := titanic-go-service
DATA_IMAGE_NAME := titanic-data-container

# Image tag. Default to v1.0.0, can be overridden, e.g., `make build TAG=v1.0.1`
TAG ?= v1.0.0

# Helm and Kind configuration
HELM_RELEASE_NAME := titanic-release
KIND_CLUSTER_NAME := titanic-cluster

DATA_SOURCE_HELM_VALUE := csv # Default data source, can be overridden

# ==============================================================================
# Targets
# ==============================================================================

# Use .PHONY to ensure these targets run even if files with the same name exist.
.PHONY: all build push svc-image data-image setup-kind install uninstall clean help

# Default target runs when you just type `make`.
all: help

## --------------------------------------
## Image Building & Pushing
## --------------------------------------

# Build both images.
build: svc-image data-image

# Install Swagger CLI tool for generating API documentation.
install-swagger:
	@echo "--> Installing Swagger"
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation for the API.
generate-swagger: install-swagger
	@echo "--> Generating Swagger documentation"
	@swag init -g cmd/server/main.go

# Seed the SQLite database with data from CSV file.
seed-sqlite:
	@echo "--> Seeding SQLite database with data from CSV"
	@go run cmd/seed/main.go

# Build the main Go service image.
svc-image: generate-swagger
	@echo "--> Building service image: $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG)"
	@docker build -f Dockerfile -t $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG) .

# Build the data container image.
data-image: seed-sqlite
	@echo "--> Building data image: $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG)"
	@docker build -f Dockerfile.data -t $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG) .

# Push both images to the configured registry.
push:
	@echo "--> Pushing images to registry: $(REGISTRY)"
	@docker push $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG)
	@docker push $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG)


## --------------------------------------
## Kubernetes Deployment (Kind & Helm)
## --------------------------------------

# Set up the local Kind cluster. This is idempotent; it won't fail if the cluster already exists.
setup-kind: build
	@if ! kind get clusters | grep -q "^$(KIND_CLUSTER_NAME)$$"; then \
		echo "--> Creating Kind cluster: $(KIND_CLUSTER_NAME)"; \
		kind create cluster --name $(KIND_CLUSTER_NAME); \
	else \
		echo "--> Kind cluster '$(KIND_CLUSTER_NAME)' already exists."; \
	fi
	@echo "--> Loading images into Kind cluster..."
	@kind load docker-image $(REGISTRY)/$(SVC_IMAGE_NAME):$(TAG) --name $(KIND_CLUSTER_NAME)
	@kind load docker-image $(REGISTRY)/$(DATA_IMAGE_NAME):$(TAG) --name $(KIND_CLUSTER_NAME)


# Install or upgrade the Helm chart.
# This passes all necessary variables to Helm, so you don't need to edit values.yaml every time.
install-helm:
	@echo "--> Installing/upgrading Helm release: $(HELM_RELEASE_NAME)"
	@helm upgrade --install $(HELM_RELEASE_NAME) ./helm/titanic-chart \
		--set image.repository=$(REGISTRY)/$(SVC_IMAGE_NAME) \
		--set image.tag=$(TAG) \
		--set dataLoaderImage.repository=$(REGISTRY)/$(DATA_IMAGE_NAME) \
		--set dataLoaderImage.tag=$(TAG) \
		--set config.dataSource=$(DATA_SOURCE_HELM_VALUE)

install: setup-kind install-helm


# Uninstall the Helm release.
uninstall:
	@echo "--> Uninstalling Helm release: $(HELM_RELEASE_NAME)"
	@helm uninstall $(HELM_RELEASE_NAME)


## --------------------------------------
## Cleanup & Help
## --------------------------------------

# Delete the local Kind cluster.
clean:
	@echo "--> Deleting Kind cluster: $(KIND_CLUSTER_NAME)"
	@kind delete cluster --name $(KIND_CLUSTER_NAME)

# Display help message.
help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build both the service and data docker images."
	@echo "  push          Push docker images to the configured registry."
	@echo "  setup-kind    Create a local Kind cluster and load images into it."
	@echo "  install       Deploy the application to the Kind cluster using Helm."
	@echo "  uninstall     Remove the application from the Kind cluster."
	@echo "  clean         Delete the local Kind cluster."
	@echo "  help          Show this help message."
	@echo ""
	@echo "You can override variables, e.g., 'make install TAG=v1.0.2 REGISTRY=myrepo'"
	@echo ""
