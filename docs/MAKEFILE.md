# Makefile Quick Reference

Quick guide to common Makefile commands for IP Verifier.

## 🚀 Most Used Commands

```bash
make help              # Show all available commands
make k8s-deploy        # Deploy everything to Kubernetes
make k8s-status        # Check deployment status
make test-api          # Test API endpoints
make k8s-logs          # View application logs
make k8s-check-updates # Check database update history
```

---

## 📋 Development

```bash
make test              # Run unit tests
make test-coverage     # Generate coverage report
make test-e2e          # End-to-end tests
make build             # Build binary
make run               # Run locally
make clean             # Remove build artifacts
make dev               # Build + run
```

---

## 🐳 Docker

```bash
make docker-build      # Build Docker image
make docker-run        # Run in container (port 8080)
make docker-stop       # Stop container
```

---

## ☸️ Kubernetes

### Deployment
```bash
make k8s-deploy        # Deploy all resources
make k8s-delete        # Delete namespace
make deploy-full       # Build Docker image + deploy
```

### Monitoring
```bash
make k8s-status        # Show pods, services, deployments, cronjobs
make k8s-logs          # Stream application logs
make k8s-check-updates # View update history and schedule
```

### Management
```bash
make k8s-restart              # Rolling restart
make k8s-scale REPLICAS=5     # Scale to N replicas
make k8s-update-db            # Manually trigger database update
```

---

## 🧪 Testing

```bash
make test              # Unit tests
make test-coverage     # Coverage report (opens coverage.html)
make test-e2e          # End-to-end tests
make test-api          # Test health and IP verification endpoints
```

---

## 💡 Common Workflows

### Initial Setup
```bash
# 1. Create .env with MaxMind credentials
./scripts/create-secret.sh

# 2. Deploy
make k8s-deploy

# 3. Verify
make k8s-status
make test-api
```

### Development Cycle
```bash
# 1. Make code changes
# 2. Test locally
make test

# 3. Build and deploy
make deploy-full

# 4. Verify
make test-api
```

### Troubleshooting
```bash
make k8s-status        # Check pod status
make k8s-logs          # View application logs
```

### Database Updates
```bash
make k8s-check-updates # Check status and schedule
make k8s-update-db     # Trigger manual update
make k8s-restart       # Restart to reload database
```

---

## 🎯 Examples

### Scale for High Traffic
```bash
make k8s-scale REPLICAS=10  # Scale up
watch make k8s-status       # Monitor
make k8s-scale REPLICAS=2   # Scale down
```

### Test Custom IPs
```bash
# Use test-api for quick checks, or:
curl -X POST http://localhost/api/v1/ip-verifier \
  -H "Content-Type: application/json" \
  -d '{"ip":"1.1.1.1","allowed_countries":["US","AU"]}' | jq .
```

---

## 🔧 Requirements

- **Docker**: For containerization
- **kubectl**: For Kubernetes operations  
- **jq**: For JSON formatting

```bash
# macOS
brew install kubectl jq

# Linux
# kubectl: https://kubernetes.io/docs/tasks/tools/
# jq: apt install jq
```

---

**See full list**: `make help`
