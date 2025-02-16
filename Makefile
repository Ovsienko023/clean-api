BINARY_NAME := app
BUILD_DIR := bin

.PHONY: all build run test clean

all: build

build:
	@echo "Build app..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

run: build
	@echo "Run app..."
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Run tests..."
	go test -v ./...

clean:
	@echo "Clean..."
	rm -rf $(BUILD_DIR)
