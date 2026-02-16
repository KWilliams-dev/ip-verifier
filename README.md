# IP Verifier

Production-ready microservice for IP address geolocation verification. Validates IP addresses against country allow-lists using MaxMind GeoIP2 databases with automated weekly updates.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-brightgreen.svg)](./Dockerfile)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5.svg)](./k8s/)

---

## 🚀 Features

- **IP Geolocation**: Identify country of origin for any IPv4/IPv6 address
- **Country Allow-Lists**: Verify IPs against configurable country codes
- **Production Ready**: Structured logging, error handling, health checks
- **Kubernetes Native**: Full K8s deployment with automated database updates
- **Minimal Footprint**: 46MB Docker image using distroless runtime
- **Auto-Updates**: Weekly GeoIP database updates via CronJob
- **Observable**: Structured JSON logging with request tracing
- **Well-Tested**: Comprehensive unit and E2E test coverage

---

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Architecture](#architecture)
- [Deployment](#deployment)
- [Configuration](#configuration)
- [Testing](#testing)
- [Documentation](#documentation)
- [Development](#development)

---

## ⚡ Quick Start

### Prerequisites

- Docker
- Kubernetes cluster (Docker Desktop, minikube, etc.)
- kubectl configured
- MaxMind GeoLite2 account ([free signup](https://www.maxmind.com/en/geolite2/signup))

### 1. Deploy to Kubernetes

```bash
# Clone repository
git clone https://github.com/KWilliams-dev/ip-verifier.git
cd ip-verifier

# Set up MaxMind credentials
echo 'ACCOUNT_ID=your_account_id' > .env
echo 'LICENSE_KEY=your_license_key' >> .env
./scripts/create-secret.sh

# Deploy all resources
make k8s-deploy

# Verify deployment
make k8s-status
```

### 2. Test the API

```bash
# Health check
curl http://localhost/api/v1/health

# Verify IP address
curl -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "allowed_countries": ["US"]
  }'
```

**Response:**
```json
{
  "ip": "8.8.8.8",
  "country_code": "US",
  "is_allowed": true
}
```

---

## 📡 API Reference

### Health Check

**Endpoint:** `GET /api/v1/health`

**Response:**
```json
{
  "status": "healthy"
}
```

### IP Verification

**Endpoint:** `POST /api/v1/ip-verifier`

**Request:**
```json
{
  "ip": "1.1.1.1",
  "allowed_countries": ["US", "CA", "AU"]
}
```

**Response (Allowed):**
```json
{
  "ip": "1.1.1.1",
  "country_code": "AU",
  "is_allowed": true
}
```

**Response (Not Allowed):**
```json
{
  "ip": "77.88.8.8",
  "country_code": "RU",
  "is_allowed": false
}
```

**Error Response:**
```json
{
  "error": "invalid IP address format"
}
```

### Status Codes

- `200 OK` - Request successful
- `400 Bad Request` - Invalid input (malformed IP, missing fields)
- `500 Internal Server Error` - Database or server error
- `503 Service Unavailable` - GeoIP database not loaded

---

## 🏗️ Architecture

### Technology Stack

- **Language**: Go 1.25
- **Framework**: Gin (HTTP router)
- **Database**: MaxMind GeoLite2-Country (MMDB format)
- **Containerization**: Docker (multi-stage build)
- **Orchestration**: Kubernetes
- **Logging**: Structured JSON logs (slog)
- **Testing**: Go test, testify

### Project Structure

```
ip-verifier/
├── cmd/
│   └── ip-verifier-api/      # Application entry point
├── internal/
│   ├── api/
│   │   └── handler/           # HTTP handlers
│   ├── config/                # Configuration management
│   ├── domain/                # Business domain interfaces
│   ├── errors/                # Custom error types
│   ├── repo/                  # GeoIP database repository
│   └── service/               # Business logic
├── k8s/                       # Kubernetes manifests
├── scripts/                   # Deployment scripts
├── docs/                      # Documentation
└── test/                      # E2E tests
```

### Design Principles

- **Clean Architecture**: Clear separation of concerns (handlers, services, repositories)
- **Dependency Injection**: Interfaces enable testability and flexibility
- **Error Handling**: Custom error types with structured error responses
- **Observability**: Structured logging with context propagation
- **Configuration**: Environment-based with sensible defaults

---

## 🚢 Deployment

### Docker

Build and run locally:

```bash
# Build image
make docker-build

# Run container
make docker-run

# Test
curl http://localhost:8080/api/v1/health
```

### Kubernetes

Full production deployment with automated database updates:

```bash
# Deploy everything
make k8s-deploy

# Check status
make k8s-status

# View logs
make k8s-logs

# Scale replicas
make k8s-scale REPLICAS=5
```

**What gets deployed:**
- Namespace: `ip-verifier`
- Deployment: 2 replicas with init container
- Service: LoadBalancer (localhost on Docker Desktop)
- PersistentVolumeClaim: 100Mi for GeoIP database
- CronJob: Weekly database updates (Wednesday 3 AM UTC)
- ConfigMap: geoipupdate configuration
- Secret: MaxMind credentials

### Database Updates

GeoIP database automatically updates weekly:

```bash
# Check update status
make k8s-check-updates

# Trigger manual update
make k8s-update-db

# View update logs
kubectl logs -n ip-verifier job/<job-name>
```

See [TESTING_GUIDE.md](./docs/TESTING_GUIDE.md) for detailed verification steps.

---

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `ENVIRONMENT` | Environment name (dev/production) | `development` |
| `GEOIP_DB_PATH` | Path to MMDB file | `data/GeoLite2-Country.mmdb` |
| `ACCOUNT_ID` | MaxMind account ID | Required for updates |
| `LICENSE_KEY` | MaxMind license key | Required for updates |

### Kubernetes ConfigMap

```yaml
GEOIPUPDATE_EDITION_IDS: GeoLite2-Country
GEOIPUPDATE_FREQUENCY: 0  # Run once (for CronJob)
GEOIPUPDATE_VERBOSE: 1
```

### Update Schedule

Modify CronJob schedule in `k8s/geoip-cronjob.yaml`:

```yaml
schedule: "0 3 * * 3"  # Every Wednesday at 3 AM UTC
```

---

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
make test

# With coverage
make test-coverage
open coverage.html
```

### E2E Tests

```bash
# Requires running service
make k8s-deploy
make test-e2e
```

### API Testing

```bash
# Quick smoke test
make test-api

# Manual tests
curl -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{"ip":"8.8.8.8","allowed_countries":["US"]}' | jq .
```

### Test Coverage

Current coverage: **90%+** across all packages

See [TESTING_GUIDE.md](./docs/TESTING_GUIDE.md) for comprehensive testing documentation.

---

## 📚 Documentation

- **[DEPLOYMENT_SUMMARY.md](./docs/DEPLOYMENT_SUMMARY.md)** - Complete deployment guide
- **[TESTING_GUIDE.md](./docs/TESTING_GUIDE.md)** - Testing procedures and verification
- **[MAKEFILE.md](./docs/MAKEFILE.md)** - Makefile command reference
- **[geoip-update.md](./docs/geoip-update.md)** - Database update mechanism details

---

## 💻 Development

### Local Development

```bash
# Install dependencies
go mod download

# Run locally (requires MMDB file)
make run

# Build binary
make build
./bin/ip-verifier-api
```

### Local Testing Without Kubernetes

Run the service locally for development:

```bash
# 1. Ensure you have the GeoIP database
# Download from MaxMind or copy from running pod:
kubectl cp ip-verifier/<pod-name>:/var/lib/geoip/GeoLite2-Country.mmdb ./data/GeoLite2-Country.mmdb -n ip-verifier

# 2. Set environment variables (optional)
export PORT=8080
export ENVIRONMENT=development
export GEOIP_DB_PATH=data/GeoLite2-Country.mmdb

# 3. Run the service
make run
# Or directly:
go run cmd/ip-verifier-api/main.go
```

**Test local service:**

```bash
# Health check
curl http://localhost:8080/api/v1/health | jq .

# Test US IP
curl -X POST http://localhost:8080/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{"ip":"8.8.8.8","allowed_countries":["US"]}' | jq .

# Test non-allowed IP
curl -X POST http://localhost:8080/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{"ip":"77.88.8.8","allowed_countries":["US"]}' | jq .
```

### Development Workflow

```bash
# 1. Make code changes
vim internal/service/ip_verifier_service.go

# 2. Run tests
make test

# 3. Build and deploy
make deploy-full

# 4. Verify
make test-api
```

### Makefile Commands

```bash
make help              # Show all commands
make test              # Run unit tests
make build             # Build binary
make docker-build      # Build Docker image
make k8s-deploy        # Deploy to Kubernetes
make k8s-status        # Check deployment status
make k8s-logs          # Stream logs
make k8s-restart       # Rolling restart
```

See [MAKEFILE.md](./docs/MAKEFILE.md) for complete command reference.

---

## 🔧 Troubleshooting

### Service Not Responding

```bash
make k8s-status        # Check pod status
make k8s-logs          # View logs
```

### Database Update Failures

```bash
# Check update jobs
kubectl get jobs -n ip-verifier -l cronjob-name=geoip-updater

# View logs
kubectl logs -n ip-verifier job/<job-name>

# Verify credentials
kubectl get secret maxmind-credentials -n ip-verifier
```

### IP Lookups Returning Errors

```bash
# Check database file exists
kubectl get pvc -n ip-verifier

# Trigger manual update
make k8s-update-db

# Restart pods
make k8s-restart
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Add unit tests for new functionality
- Update documentation for API changes
- Use structured logging (slog)
- Handle errors explicitly

---

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 🙏 Acknowledgments

- [MaxMind](https://www.maxmind.com) - GeoLite2 database
- [oschwald/geoip2-golang](https://github.com/oschwald/geoip2-golang) - GeoIP2 Go library
- [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP framework

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/KWilliams-dev/ip-verifier/issues)
- **Documentation**: [./docs](./docs)
- **MaxMind Support**: [MaxMind Support](https://support.maxmind.com)

---

**Built with ❤️ using Go and Kubernetes**
