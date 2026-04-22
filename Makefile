BIN_DIR := bin
BIN_NAME := csv-decode

.PHONY: build run test clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) .

run: build
	./$(BIN_DIR)/$(BIN_NAME)

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)
