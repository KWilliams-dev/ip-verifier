.PHONY: help test test-coverage test-e2e build run clean docker-build docker-run docker-stop k8s-deploy k8s-delete k8s-status k8s-logs k8s-restart k8s-update-db k8s-check-updates test-api

APP_NAME=ip-verifier
NAMESPACE=ip-verifier
DOCKER_IMAGE=$(APP_NAME):local
PORT=8080

##@ Help

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

test: ## Run all tests
	@go test ./... -v

test-coverage: ## Generate coverage report
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-e2e: ## Run E2E tests
	@go test -tags=e2e ./test/e2e/ -v

build: ## Build binary
	@CGO_ENABLED=0 go build -o bin/$(APP_NAME)-api cmd/ip-verifier-api/main.go

run: ## Run locally
	@go run cmd/ip-verifier-api/main.go

clean: ## Clean artifacts
	@rm -rf bin/ coverage.out coverage.html

##@ Docker

docker-build: ## Build image
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-stop ## Run container
	@docker run -d -p $(PORT):8080 \
		-e PORT=8080 \
		-e GEOIP_DB_PATH=/app/data/GeoLite2-Country.mmdb \
		--name $(APP_NAME) $(DOCKER_IMAGE)
	@sleep 2 && curl -s http://localhost:$(PORT)/api/v1/health | jq .

docker-stop: ## Stop container
	@docker stop $(APP_NAME) 2>/dev/null || true
	@docker rm $(APP_NAME) 2>/dev/null || true

##@ Kubernetes

k8s-deploy: ## Deploy all resources
	@kubectl apply -f k8s/namespace.yaml
	@./scripts/create-secret.sh 2>/dev/null || true
	@kubectl apply -f k8s/
	@make k8s-status

k8s-delete: ## Delete namespace
	@kubectl delete namespace $(NAMESPACE)

k8s-status: ## Show status
	@kubectl get pods,svc,deploy,cronjobs -n $(NAMESPACE)

k8s-logs: ## Stream logs
	@kubectl logs -n $(NAMESPACE) -l app=$(APP_NAME) --tail=50 -f

k8s-restart: ## Rolling restart
	@kubectl rollout restart deployment/$(APP_NAME) -n $(NAMESPACE)
	@kubectl rollout status deployment/$(APP_NAME) -n $(NAMESPACE)

k8s-scale: ## Scale (REPLICAS=N)
	@kubectl scale deployment/$(APP_NAME) -n $(NAMESPACE) --replicas=$(REPLICAS)

k8s-update-db: ## Trigger DB update
	@kubectl create job -n $(NAMESPACE) geoip-manual-$$(date +%s) --from=cronjob/geoip-updater

k8s-check-updates: ## Show update history
	@kubectl get jobs -n $(NAMESPACE) -l cronjob-name=geoip-updater --sort-by=.status.startTime | tail -6
	@kubectl get cronjob geoip-updater -n $(NAMESPACE)

##@ Testing

test-api: ## Test endpoints
	@curl -s http://localhost/api/v1/health | jq .
	@curl -s -X POST http://localhost/api/v1/ip-verifier \
		-H "Content-Type: application/json" \
		-d '{"ip":"8.8.8.8","allowed_countries":["US"]}' | jq .

##@ Workflows

dev: build run ## Build and run

deploy-full: docker-build k8s-deploy ## Build + deploy

.DEFAULT_GOAL := help

