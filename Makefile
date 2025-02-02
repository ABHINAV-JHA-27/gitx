build:
	@echo "Building the application..."
	@go mod tidy
	@go build -o app cmd/main.go
	@echo "Application built successfully!!"