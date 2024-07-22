


test:
	@echo "Running tests..."
	@echo "Testing ./decoders"
	@go test ./decoders
	@echo "Testing ./encoders"
	@go test ./encoders
	@echo "Ended tests..."
