
client:
	@go run ./cmd/client/main.go

server:
	@go run ./cmd/server/main.go

test:
	@echo Running tests...
	@echo Testing ./decoders
	@go test -v ./decoders
	@echo Testing ./encoders
	@go test -v ./encoders
	@echo Ended tests...
