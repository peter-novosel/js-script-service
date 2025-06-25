APP_NAME=js-script-service
LAMBDA_MAIN=./cmd/lambda
SERVER_MAIN=./cmd/server
BIN_DIR=./bin

run: ## Run server locally
	go run $(SERVER_MAIN)

lint: ## Run linter
	golangci-lint run

test: ## Run all tests
	go test ./...

build: ## Build binary for local use
	go build -o $(BIN_DIR)/server $(SERVER_MAIN)

build-lambda: ## Compile for AWS Lambda (linux/amd64)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/bootstrap $(LAMBDA_MAIN)

deploy-lambda: build-lambda ## Zip and prepare for Lambda deployment
	cd $(BIN_DIR) && zip lambda.zip bootstrap

clean:
	rm -rf $(BIN_DIR)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
