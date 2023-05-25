# Go compiler
GO := go

# Name of your binary
BINARY_NAME := rebaseAllDirs

# Directories
SRC_DIR := ./
BUILD_DIR := ./build
BIN_DIR := ./bin

# Source files
SRCS := $(shell find $(SRC_DIR) -name '*.go')

# TODO: update this in your local testing
# Variables
REPOSROOT := "/Users/nawazhussain/msftcode"
DEFAULTUPSTREAM := "main"

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

test: clean all
	@echo "Running \"${BINARY_NAME}\""
	@./${BIN_DIR}/${BINARY_NAME} --reposRoot $(REPOSROOT) --defaultUpstream $(DEFAULTUPSTREAM)
