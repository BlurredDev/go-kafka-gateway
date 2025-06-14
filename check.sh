# 1. Tidy dependencies
go mod tidy

# 2. Build all packages
go build ./...

# 3. Run tests with coverage, excluding stress package
go test -v -coverpkg=./... -coverprofile=coverage.out $(go list ./... | grep -v '/internal/stress') -covermode=atomic
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 4. Static analysis
go vet ./...

# 5. Format check (non-destructive)
gofmt -l .
