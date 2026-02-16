# Go Commands Reference

## Running the Application
```bash
# Run the application
go run cmd/retrovault-api/main.go

# Build the application
go build -o retrovault cmd/retrovault-api/main.go

# Run the built binary
./retrovault
```

## Dependency Management
```bash
# Initialize go module
go mod init github.com/KWilliams-dev/retro-vault

# Add a dependency
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/lib/pq

# Download all dependencies
go mod download

# Tidy up dependencies (remove unused)
go mod tidy

# Verify dependencies
go mod verify
```

## Code Formatting & Quality
```bash
# Format a specific file
go fmt internal/service/game_service.go
go fmt internal/domain/game.go

# Format entire project
go fmt ./...

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build without running
go build ./...
```