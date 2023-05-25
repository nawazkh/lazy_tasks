# Go compiler
GO := go

# Name of your binary
BINARY_NAME := cleanupAzureResources

# Directories
SRC_DIR := ./
BUILD_DIR := ./build
BIN_DIR := ./bin

# Source files
SRCS := $(shell find $(SRC_DIR) -name '*.go')

# TODO: update this in your local testing
# Variables
AZURE_SUB := $(shell echo $$AZURE_SUBSCRIPTION_ID)

.PHONY: all build clean test

all: build

build: $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME): $(SRCS)
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)
	@mv $(BUILD_DIR)/$(BINARY_NAME) $(BIN_DIR)
	@echo "Build complete."

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(BIN_DIR)
	@echo "Clean complete."

check-env:
    ifndef AZURE_SUBSCRIPTION_ID
        $(error AZURE_SUBSCRIPTION_ID is not set in the environment)
    endif
    ifndef AZURE_CLIENT_ID
        $(error AZURE_CLIENT_ID is not set in the environment)
    endif
    ifndef AZURE_CLIENT_SECRET
        $(error AZURE_CLIENT_SECRET is not set in the environment)
    endif
	ifndef AZURE_TENANT_ID
        $(error AZURE_TENANT_ID is not set in the environment)
    endif

.PHONY: check-env

test: check-env clean all 
	@echo "Running \"${BINARY_NAME}\""
	@./${BIN_DIR}/${BINARY_NAME} --reposRoot $(REPOSROOT)
