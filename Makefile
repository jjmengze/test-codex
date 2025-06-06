# Makefile for Go project

BINARY_NAME := log-receiver
PKG := ./cmd/server
OUTPUT_DIR := bin

# Build target platform
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
PORT = 8080
IS_TEST_PEM ?= true
.PHONY: all build build-all test run fmt docker-build docker-run

all: build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH) $(PKG)

# Build for common platforms
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 windows/amd64
OUTPUT_DIR := build

build-all:
	@mkdir -p $(OUTPUT_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%%/*}; arch=$${platform##*/}; \
		output=$(OUTPUT_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		[ $$os = "windows" ] && output=$$output.exe; \
		echo "building for $$os/$$arch"; \
		GOOS=$$os GOARCH=$$arch go build -o $$output $(PKG); \
	done

run:
	go run $(PKG) --port $(PORT) --is_test_pem $(IS_TEST_PEM)

test:
	go test -v ./...

# Coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Lint 檢查 (需安裝 golangci-lint)
lint: tool
	golangci-lint run ./...

# 自動產生 mock（使用 mockery v2）
mock: tool
	mockery --all --keeptree --output=./mock --outpkg=mock


generate: tool
	go generate ./...

tool:
	# Lint 工具
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

    # Mockery
	go install github.com/vektra/mockery/v2@latest

    # Reflex (即時測試)
	go install github.com/cespare/reflex@latest

fmt:
	go fmt ./...

# Docker image config
IMAGE_NAME ?= log-receiver
IMAGE_TAG ?= latest
IMAGE := $(IMAGE_NAME):$(IMAGE_TAG)
PLATFORMS := linux/amd64,linux/arm64

docker-build:
	docker buildx build  --tag $(IMAGE) --output type=docker .

docker-run:
	docker run --rm \
		--env PORT=$(PORT) \
		--env IS_TEST_PEM=$(IS_TEST_PEM) \
		--env-file .env \
		-p 8080:8080 \
		$(IMAGE)

# Optional: clean local image
docker-clean:
	docker rmi $(IMAGE) || true

# terraform block (IAC)

# Terraform workdir
TF_DIR := terraform
TF_VARS := terraform.tfvars


.PHONY: tf-init
tf-init:
	cd $(TF_DIR) && terraform init

.PHONY: tf-plan
tf-plan: tf-init
	cd $(TF_DIR) && terraform plan -var-file=$(TF_VARS) -out=tfplan


.PHONY: tf-apply
tf-apply: tf-plan
	@read -p "⚠️  Are you sure you want to apply this Terraform plan? (y/N) " ans && [ "$$ans" = "y" ] && \
	cd $(TF_DIR) && terraform apply tfplan || echo "❌ Cancelled."

.PHONY: tf-destroy
tf-destroy: tf-init
	@read -p "⚠️  WARNING: This will destroy all resources. Continue? (y/N) " ans && [ "$$ans" = "y" ] && \
	cd $(TF_DIR) && terraform destroy -var-file=$(TF_VARS) || echo "❌ Cancelled."

.PHONY: clean
clean: tf-destroy
	rm -rf $(OUTPUT_DIR)