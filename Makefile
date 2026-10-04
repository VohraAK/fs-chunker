APP_NAME=fs-chunker
BIN_DIR=bin

.PHONY: all build run clean test
all: build
	
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BIN_DIR)/$(APP_NAME) main.go
run:
	go run main.go $(ARGS)
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
test:
	go test ./...