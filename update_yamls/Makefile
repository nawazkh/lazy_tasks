# Name of the binary
BINARY_NAME = update_yaml

# Go package name
PACKAGE_NAME = .

# Input YAML file
INPUT_YAML ?= input.yaml

# Default target: build the binary
all: build

# Build the binary
build:
	go build -o $(BINARY_NAME) $(PACKAGE_NAME)

# Install the binary to GOPATH/bin
install:
	go install $(PACKAGE_NAME)

# Run the binary
run: clean build
	./$(BINARY_NAME) $(INPUT_YAML)

# Clean up the build files
clean:
	rm -f $(BINARY_NAME)

# Print help
help:
	@echo "Makefile commands:"
	@echo "  all       - Build the binary (default)"
	@echo "  build     - Build the binary"
	@echo "  install   - Install the binary to GOPATH/bin"
	@echo "  run       - Run the binary"
	@echo "  clean     - Clean up build files"
	@echo "  help      - Print this help message"
