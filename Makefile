all: test

generate:
	go run cmd/generate-routers/main.go && gofumpt -l -w .

test: generate
	go test ./...